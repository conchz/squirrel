package app

func init() {
	OnAppStart(InitMySQL)
	OnAppStart(InitRedis)
}

func InitMySQL() {

}

func InitRedis() {

}
