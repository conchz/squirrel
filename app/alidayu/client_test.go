package alidayu_test

import (
	"fmt"
	"github.com/lavenderx/squirrel/app/alidayu"
	"io/ioutil"
	"os"
)

// Example of Client.Post()
func ExampleClient_Post() {
	// Create a new client(创建一个新的Client实例).
	c := &alidayu.Client{AppKey: "", AppSecret: "", UseHTTPS: false}

	// ---------------------------------------
	// Send Verification Code in SMS(发送短信验证码).
	// ---------------------------------------
	// Set parameters(设置API所需的所有参数，包括公共参数和方法相关参数).
	params := map[string]string{}

	// Common parameters will be filled with default values automatically.
	// You may also set them to override default ones.
	// 公共参数会自动使用默认值填充，不需要用户设置.
	// 当然，用户可以使用自己的设置来覆盖默认值.
	//   params["format"] = "json"      // "json" or "xml"
	//   params["v"] = "2.0"            // "2.0" by default
	//   params["sign_method"] = "md5"  // "md5" or "hmac"

	// No need to set "timestamp" and "sign", these parameters will be calculated and filled automatically.
	// 不需要用户设置"timestamp"和"sign",这些参数会被自动计算和设置.

	// Set method specified parameters(设置方法相关特定参数).
	params["method"] = "alibaba.aliqin.fc.sms.num.send"           // Set method to send SMS(API接口名称).
	params["sms_type"] = "normal"                                 // Set SMS type(短信类型).
	params["sms_free_sign_name"] = ""                             // Set SMS signature(短信签名).
	params["sms_param"] = `{"code":"123456", "product":"My App"}` // Set variable for SMS template(短信模板变量).
	params["sms_template_code"] = ""                              // Set SMS template code(短信模板ID).
	params["rec_num"] = ""                                        // Set phone number to send SMS(短信接收号码).

	// Call Post() to post the request.
	resp, err := c.Post(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "c.Post() error: %v\n", err)
		return
	}

	// Read HTTP Response.
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ioutil.ReadAll() error:%v\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "c.Post() successfully\n%v\n", string(data))

	// ------------------------------------------
	// Send Verification Code in Single Call(发送文本转语音通知验证码).
	// ------------------------------------------
	// Set Parameters.
	params = map[string]string{}
	params["method"] = "alibaba.aliqin.fc.tts.num.singlecall"     // Set method to make single call(API接口名称).
	params["tts_param"] = `{"code":"123456", "product":"My App"}` // Set variable for TTS template(文本转语音（TTS）模板变量).
	params["called_num"] = ""                                     // Set phone number to make single call(被叫号码).
	params["called_show_num"] = ""                                // Set show number(被叫号显).
	params["tts_code"] = ""                                       // Set TTS code(TTS模板ID).

	// Call Post() to post the request.
	resp2, err := c.Post(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "c.Post() error: %v\n", err)
		return
	}

	// Read HTTP Response.
	defer resp2.Body.Close()
	data2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ioutil.ReadAll() error:%v\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "c.Post() successfully\n%v\n", string(data2))

	// Output:
}

// Example of Client.Exec()
func ExampleClient_Exec() {
	// Create a new client(创建一个新的Client实例).
	c := &alidayu.Client{AppKey: "", AppSecret: "", UseHTTPS: false}

	// ---------------------------------------
	// Send Verification Code in SMS(发送短信验证码).
	// ---------------------------------------
	// Set parameters(设置API所需的所有参数，包括公共参数和方法相关参数).
	params := map[string]string{}

	// Common parameters will be filled with default values automatically.
	// You may also set them to override default ones.
	// 公共参数会自动使用默认值填充，不需要用户设置.
	// 当然，用户可以使用自己的设置来覆盖默认值.
	//   params["format"] = "json"      // "json" or "xml"
	//   params["v"] = "2.0"            // "2.0" by default
	//   params["sign_method"] = "md5"  // "md5" or "hmac"

	// No need to set "timestamp" and "sign", these parameters will be calculated and filled automatically.
	// 不需要用户设置"timestamp"和"sign",这些参数会被自动计算和设置.

	// Set method specified parameters(设置方法相关特定参数).
	params["method"] = "alibaba.aliqin.fc.sms.num.send"           // Set method to send SMS(API接口名称).
	params["sms_type"] = "normal"                                 // Set SMS type(短信类型).
	params["sms_free_sign_name"] = ""                             // Set SMS signature(短信签名).
	params["sms_param"] = `{"code":"123456", "product":"My App"}` // Set variable for SMS template(短信模板变量).
	params["sms_template_code"] = ""                              // Set SMS template code(短信模板ID).
	params["rec_num"] = ""                                        // Set phone number to send SMS(短信接收号码).

	// Call Exec() to post the request.
	success, result, err := c.Exec(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "c.Exec() error: %v\nsuccess: %v\nresult: %v\n", err, success, result)
		return
	}

	fmt.Fprintf(os.Stderr, "c.Exec() successfully\nsuccess: %v\nresult: %s\n", success, result)

	// ------------------------------------------
	// Send Verification Code in Single Call(发送文本转语音通知验证码).
	// ------------------------------------------
	// Set Parameters.
	params = map[string]string{}
	params["method"] = "alibaba.aliqin.fc.tts.num.singlecall"     // Set method to make single call(API接口名称).
	params["tts_param"] = `{"code":"123456", "product":"My App"}` // Set variable for TTS template(文本转语音（TTS）模板变量).
	params["called_num"] = ""                                     // Set phone number to make single call(被叫号码).
	params["called_show_num"] = ""                                // Set show number(被叫号显).
	params["tts_code"] = ""                                       // Set TTS code(TTS模板ID).

	// Call Exec() to post the request.
	success, result, err = c.Exec(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "c.Exec() error: %v\nsuccess: %v\nresult: %v\n", err, success, result)
		return
	}

	fmt.Fprintf(os.Stderr, "c.Exec() successfully\nsuccess: %v\nresult: %s\n", success, result)

	// Output:
}
