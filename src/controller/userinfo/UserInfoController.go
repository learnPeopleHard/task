package userinfo

import (
	"encoding/json"
	"fmt"
	"gotest1/src/common"
	"gotest1/src/logic/check"
	"log"
	"net/http"
	"strconv"
	"time"
)

//执行登录
func UpdateTask(resp http.ResponseWriter, req *http.Request)  {
	fmt.Println("request  start:")
	tokenObj,err:= check.CheckLoginAuthGet(resp,req)
	if err!=nil{
		fmt.Fprint(resp, err.Error())
		return
	}
	req.ParseForm()
	param_userId, found3 := req.Form["userId"]
	param_nickName, found1 := req.Form["nickName"]
	param_picUrl, found2 := req.Form["picUrl"]
	if !(found1 && found2 && found3) {
		fmt.Fprint(resp, "nickName picUrl userId不能为空")
		return
	}
	nickName := param_nickName[0]
	picUrl := param_picUrl[0]
	if nickName=="" || picUrl=="" {
		fmt.Fprint(resp, "nickName picUrl不能为空")
		return
	}
	userId,err := strconv.ParseInt(param_userId[0], 10, 64)
	if err!=nil{
		fmt.Fprint(resp, "用户ID非法")
		return
	}
	s := "userName: "+tokenObj.UserName+"nickName:" + nickName + ",picUrl:" + picUrl
	fmt.Println(s)
	loginResponseDTO:= rpcUpdate(userId,nickName,picUrl,tokenObj.UserName)
	//向客户端发送json数据
	bytes, _ := json.Marshal(loginResponseDTO)
	fmt.Fprint(resp, string(bytes))
}

func QueryUserInfo(resp http.ResponseWriter, req *http.Request)  {
	startTime:=time.Now().UnixMilli()
	fmt.Printf("====request time is: %d \n",startTime)
	tokenObj,err:= check.CheckLoginAuthGet(resp,req)
	if err!=nil{
		fmt.Fprint(resp, err.Error())
		return
	}
	fmt.Printf("====userId: %d auth time is: %d \n", tokenObj.UserId,time.Now().UnixMilli()-startTime)
	startTime=time.Now().UnixMilli()
	userInfoDTO:= queryUserInfo(tokenObj.UserId)
	//向客户端发送json数据
	bytes, _ := json.Marshal(userInfoDTO)
	fmt.Printf("====userId: %d QueryUserInfo time is: %d \n", tokenObj.UserId,time.Now().UnixMilli()-startTime)
	fmt.Fprint(resp, string(bytes))
}

func queryUserInfo(userId int64) UserInfoDTO {
	fmt.Println("QueryUserInfo request  start,userId:",userId)
	conn, err := common.ClientPool.Get()
	if err != nil {
		log.Printf("连接池没有链接了 %v \n",err)
		return UserInfoDTO{}
	}
	defer common.ClientPool.Put(conn)
	var reply UserInfoDTO
	for {
		userIdStr :=strconv.FormatInt(userId,10)
		err = conn.Call("Listener.QueryUserInfo", []byte(userIdStr), &reply)
		if err != nil {
			log.Printf("查询QueryUserInfo 报错 %v \n",err)
			return UserInfoDTO{}
		}
		log.Printf("Reply: %v ", reply )
		return reply
	}
}

func rpcUpdate(userId int64,nickName string,profilePicture string,userName string) common.BaseResponseDTO{
	conn, err := common.ClientPool.Get()
	if err != nil {
		log.Fatal(err)
	}
	defer common.ClientPool.Put(conn)
	for {
		var req  = UserInfoDTO{}
		req.UserId = userId
		req.Nickname=nickName
		req.ProfilePicture = profilePicture
		req.UserName =userName
		var reply common.BaseResponseDTO
		str,err:= InfoEncode(req)
		if err != nil {
			log.Fatal(err)
			reply = common.BaseResponseDTO{}
			reply.Code=112
			reply.Message="序列化失败"
			return reply
		}
		err = conn.Call("Listener.UpdateUserInfo", str, &reply)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Reply: %v, Data: %v", reply, reply.Code)
		return reply
	}
}
