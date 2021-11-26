package login

import (
	"encoding/json"
	"fmt"
	"gotest1/src/common"
	"log"
	"net/http"
	"time"
)

//执行登录
func LoginTask(resp http.ResponseWriter, req *http.Request)  {
	fmt.Println("request  start:")
	startTime:=time.Now().UnixMilli()
	req.ParseForm()
	paramUsername, found1 := req.Form["userName"]
	paramPassword, found2 := req.Form["passWord"]
	if !(found1 && found2) {
		fmt.Fprint(resp, "用户名和密码必填")
		return
	}
	userName := paramUsername[0]
	passWord := paramPassword[0]
	if userName=="" || passWord==""{
		fmt.Fprint(resp, "用户名和密码必填")
		return
	}
	s := "userName:" + userName + ",password:" + passWord
	fmt.Println(s)
	loginResponseDTO:=rpcLogin(userName,passWord)
	if loginResponseDTO.Code==100{
		cookie := &http.Cookie{
			Name:   common.RedisToken,
			Value:  loginResponseDTO.LoginToken,
			MaxAge: 3600,
			Domain: "localhost",
			Path:   "/",
		}
		http.SetCookie(resp, cookie)
	}
	//向客户端发送json数据
	bytes, _ := json.Marshal(loginResponseDTO)
	fmt.Printf("userName %s login time is %d resp: %s \n", userName, time.Now().UnixMilli()-startTime,string(bytes))
	fmt.Fprint(resp,string(bytes))
}

func  rpcLogin(name string,password string) LoginResponseDTO {
	//client, err := rpc.Dial("tcp", "localhost:12345")
	conn,err := common.ClientPool.Get()
	if err != nil {
		fmt.Println("Get error: ", err)
		return LoginResponseDTO{Code: 122,Message: err.Error()}
	}
	defer common.ClientPool.Put(conn)
	if err != nil {
		log.Printf("create client err:%s\n", err)
		return LoginResponseDTO{Code: 122,Message: err.Error()}
	}
	for {
		var req  = LoginRequestDTO{}
		req.Name = name
		req.Password=password
		var reply LoginResponseDTO
		str,err:=reqencode(req)
		if err != nil {
			log.Fatal(err)
			reply = LoginResponseDTO{}
			reply.Code=112
			reply.Message="序列化失败"
			return reply
		}
		err = conn.Call("Listener.Login", str, &reply)
		if err != nil {
			log.Println("login fail ",name,err )
			log.Fatal(err)
		}
		//如果登录成功则记录token
		if reply.Code==100{
			redisTokenStr,err:=common.BuildRedisToken(reply.UserId,name)
			if err!=nil{
				log.Printf("redisTokenStr %s build异常 %s ",redisTokenStr, err.Error())
				return LoginResponseDTO{Code: 111,Message: "内部错误请联系XXX"}
			}
			webTokenStr,err:=common.BuildWebToken(reply.UserId,name,time.Now().UnixMilli())
			if err!=nil{
				log.Printf("webTokenStr %s build异常 %s ",webTokenStr, err.Error())
				return LoginResponseDTO{Code: 111,Message: "内部错误请联系XXX"}
			}
			var redisToken = common.TokenDTO{}
			redisToken.Token = redisTokenStr
			redisToken.UserId = reply.UserId
			tokenEncode,err := common.Encode(redisToken)
			log.Printf("name：%s 加密的token: %v 开始set token\n", name, redisToken)
			var replyToken common.BaseResponseDTO
			err = conn.Call("Listener.RedisSetToken",tokenEncode , &replyToken)
			reply = LoginResponseDTO{}
			if err!=nil{
				reply.Code=112
				reply.Message="登录失败:"+err.Error()
				return reply
			}
			if replyToken.Code!=100{
				reply.Code=112
				reply.Message="登录失败:"+replyToken.Message
				return reply
			}
			reply.Code=100
			reply.Message="登录成功"
			reply.LoginToken = webTokenStr
		}

		log.Printf("Reply: %v, Data: %v", reply, reply.Code)
		return reply
	}
}
