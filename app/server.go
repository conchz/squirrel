package app

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lavenderx/squirrel/app/crypto"
	"github.com/lavenderx/squirrel/app/model"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"time"
)

type (
	StartupHook struct {
		order int
		f     func()
	}

	ShutdownHook struct {
		order int
		f     func()
	}

	StartupHooks []StartupHook

	ShutdownHooks []ShutdownHook
)

var (
	isDebug            bool
	port               int
	logger             *zap.SugaredLogger
	startupHooks       StartupHooks
	shutdownHooks      ShutdownHooks
	staticFilesHandler echo.HandlerFunc
)

func (slice StartupHooks) Len() int {
	return len(slice)
}

func (slice StartupHooks) Less(i, j int) bool {
	return slice[i].order < slice[j].order
}

func (slice StartupHooks) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice ShutdownHooks) Len() int {
	return len(slice)
}

func (slice ShutdownHooks) Less(i, j int) bool {
	return slice[i].order < slice[j].order
}

func (slice ShutdownHooks) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func Run() {
	runStartupHooks()

	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler
	e.Logger.SetOutput(ioutil.Discard)

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("logger", logger)
			return h(c)
		}
	})
	e.Use(recoverWithConfig(middleware.DefaultRecoverConfig))
	e.Use(middleware.Secure())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderWWWAuthenticate,
			echo.HeaderAuthorization},
		AllowMethods: []string{
			echo.GET,
			echo.PUT,
			echo.POST,
			echo.DELETE},
	}))

	// Api group
	apiGroup := e.Group("/api/v1")
	// Configure middleware with the custom claims type for api group
	apiGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &JWTClaims{},
		SigningKey: []byte("secret"),
	}))
	apiGroup.GET("", api)

	// serves the index.html and other static files from rice
	e.GET("/", staticFilesHandler)
	e.GET("/static/*", staticFilesHandler)

	// Login route
	e.POST("/login", login)

	// Add shutdown hook
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			logger.Infof("Server will be closed, which is triggered by %s", sig.String())

			runShutdownHooks()

			logger.Infof("Server closed on %s", getLocalIP())
			os.Exit(1)
		}
	}()

	// Start server
	address := fmt.Sprintf(":%v", port)
	logger.Infof("Server started on [::]%v", address)
	e.Logger.Fatal(e.StartServer(&http.Server{
		Addr:         address,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}))
}

func runStartupHooks() {
	sort.Sort(startupHooks)
	for _, hook := range startupHooks {
		hook.f()
	}
}

func runShutdownHooks() {
	sort.Sort(shutdownHooks)
	for _, hook := range shutdownHooks {
		hook.f()
	}
}

// OnAppStart registers a function to be run at app startup.
//
// The order you register the functions will be the order they are run.
// You can think of it as a FIFO queue.
// This process will happen after the config file is read
// and before the server is listening for connections.
//
// Ideally, your application should have only one call to init() in the file init.go.
// The reason being that the call order of multiple init() functions in
// the same package is undefined.
// Inside of init() call OnAppStart() for each function you wish to register.
//
// This can be useful when you need to establish connections to databases or third-party services,
// setup app components, compile assets, or any thing you need to do between starting Server and accepting connections.
func OnAppStart(f func(), order ...int) {
	o := 1
	if len(order) > 0 {
		o = order[0]
	}
	startupHooks = append(startupHooks, StartupHook{order: o, f: f})
}

func OnAppStop(f func(), order ...int) {
	o := 1
	if len(order) > 0 {
		o = order[0]
	}
	shutdownHooks = append(shutdownHooks, ShutdownHook{order: o, f: f})
}

// https://jonathanmh.com/building-a-golang-api-with-echo-and-mysql/
// https://www.netlify.com/blog/2016/10/20/building-a-restful-api-in-go/
// https://xiequan.info/go%E4%B8%8Ejson-web-token/
// https://zhuanlan.zhihu.com/p/26300634

type (
	JWTClaims struct {
		UserId    int64  `json:"user_id"`
		Username  string `json:"user_name"`
		Cellphone string `json:"cellphone"`
		Email     string `json:"email"`
		jwt.StandardClaims
	}

	httpError struct {
		code    int
		Key     string `json:"error"`
		Message string `json:"message"`
	}
)

func newHTTPError(code int, key string, msg string) *httpError {
	return &httpError{
		code:    code,
		Key:     key,
		Message: msg,
	}
}

// Error makes it compatible with `error` interface.
func (e *httpError) Error() string {
	return e.Key + ": " + e.Message
}

// httpErrorHandler customize echo's HTTP error handler.
func httpErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		key  = "InternalServerError"
		msg  string
	)

	if he, ok := err.(*httpError); ok {
		code = he.code
		key = he.Key
		msg = he.Message
	} else if ehe, ok := err.(*echo.HTTPError); ok {
		code = ehe.Code
		key = http.StatusText(code)
		msg = key
	} else if isDebug {
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err := c.NoContent(code)
			if err != nil {
				logger.Error(err)
			}
		} else {
			err := c.JSON(code, newHTTPError(code, key, msg))
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

// Extend echo's middleware: https://echo.labstack.com/cookbook/middleware
//
// Customize echo's RecoverWithConfig to use zap log
func recoverWithConfig(config middleware.RecoverConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = middleware.DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					var err error
					switch r := r.(type) {
					case error:
						err = r
					default:
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						c.Get("logger").(*zap.SugaredLogger).Errorf("[%s] %s %s\n", "PANIC RECOVER", err, stack[:length])
					}

					c.Error(err)
				}
			}()

			return next(c)
		}
	}
}

// curl -i -w "\n" -H "'Content-Type': 'application/json; charset=UTF-8'" -d "username=test&password=testSecret" http://localhost:7000/login
// {
//   "token": "××××××××××××××××"
// }
func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	encryptedPasswd := crypto.EncryptPassword([]byte(password))

	if u := mySQLTemplate.GetByNonEmptyFields(&model.User{
		Username: username,
		Password: encryptedPasswd,
	}); u != nil {
		user := u.(*model.User)
		if username == user.Username && encryptedPasswd == user.Password {
			claims := &JWTClaims{
				user.Id,
				user.Username,
				user.Cellphone,
				user.Email,
				jwt.StandardClaims{
					ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
				},
			}

			// Create token with claims
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			// Generate encoded token and send it as response.
			t, err := token.SignedString([]byte("secret"))
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, echo.Map{
				"token": t,
			})
		}
	}

	return echo.ErrUnauthorized
}

func (claims JWTClaims) Valid() error {
	if err := claims.StandardClaims.Valid(); err != nil {
		vErr := err.(*jwt.ValidationError)
		return newHTTPError(http.StatusUnauthorized, "TokenValidError", vErr.Inner.Error())
	}

	if claims.UserId > 0 && claims.Username != "" {
		return nil
	}

	return newHTTPError(http.StatusUnauthorized,
		"TokenValidError",
		"Must provide user_id and user_name")
}

// curl -i -w "\n" -H "Authorization: Bearer $token" http://localhost:7000/api/v1
func api(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	username := claims.Username
	return c.JSON(http.StatusOK, "Hello, "+username+"!")
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.Errorf("Oops: %s", err.Error())
		return ""
	}

	var ipv4Addrs = []string{}

	for _, a := range addrs {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				//os.Stdout.WriteString(ipNet.IP.String() + "\n")
				ipv4Addrs = append(ipv4Addrs, ipNet.IP.String())
			}
		}
	}

	if len(ipv4Addrs) == 0 {
		return ""
	}

	return ipv4Addrs[0]
}