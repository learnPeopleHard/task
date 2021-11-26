package main

import (
	"fmt"
	"gotest1/src/common"
	"gotest1/src/controller/file"
	"gotest1/src/controller/login"
	"gotest1/src/controller/userinfo"
	"net/http"
)



func main()  {
	fmt.Println("ListenAndServe start:")
	http.HandleFunc("/login",login.LoginTask)

	http.HandleFunc("/update", userinfo.UpdateTask)

	http.HandleFunc("/queryUserInfo", userinfo.QueryUserInfo)

	http.HandleFunc("/upload", file.UploadFileHandler)
	common.ClientPoolInit()
	//服务器要监听的主机地址和端口号
	err := http.ListenAndServe("127.0.0.1:8081", nil)
	if err != nil {
		fmt.Println("ListenAndServe error:", err.Error())
	}
}

