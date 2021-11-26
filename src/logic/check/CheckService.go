package check

import (
	"errors"
	"fmt"
	"gotest1/src/common"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func CheckLoginAuthGet(w http.ResponseWriter,r *http.Request) (common.TokenDTO,error)  {
	startTime:=time.Now().UnixMilli()
	r.ParseForm()
	paramUserid, found3 := r.Form["userId"]
	if !(found3) {
		return common.TokenDTO{},errors.New("userId不能为空")
	}
	redisToken, found4 := r.Form["redisToken"]
	var webToken string
	if !(found4){
		cReq, err := r.Cookie(common.RedisToken)
		if err!=nil{
			return  common.TokenDTO{},errors.New("token失效，请登录")
		}
		fmt.Printf("cookie:%#v, err:%v\n", cReq, err)
		if cReq==nil{
			return  common.TokenDTO{},errors.New("token失效，请登录")
		}
	}else{
		webToken = redisToken[0]
	}
	userId,err := strconv.ParseInt(paramUserid[0], 10, 64)
	if err!=nil{
		return common.TokenDTO{},errors.New("userId非法"+err.Error())
	}
	webTokenObj,err := auth(webToken,userId)
	if err!=nil{
		return common.TokenDTO{},err
	}
	if webToken!=webTokenObj.Token{
		cResp := http.Cookie{
			Name: common.RedisToken,
			Value:webTokenObj.Token,
		}
		http.SetCookie(w,&cResp)
	}
	fmt.Printf("auth time is: %d  \n", time.Now().UnixMilli()-startTime)
	return webTokenObj,nil
}

func CheckLoginAuthPost(w http.ResponseWriter,r *http.Request) (common.TokenDTO,error)  {
	userId,err := strconv.ParseInt(r.PostFormValue("userId"),10,64)
	if err!=nil{
		return common.TokenDTO{},errors.New("userId非法")
	}
	redisToken := r.PostFormValue("redisToken")
	var webToken string
	if redisToken==""{
		cReq, err := r.Cookie(common.RedisToken)
		if err!=nil{
			return  common.TokenDTO{},errors.New("token失效，请登录")
		}
		fmt.Printf("cookie:%#v, err:%v\n", cReq, err)
		if cReq==nil{
			return  common.TokenDTO{},errors.New("token失效，请登录")
		}
	}else{
		webToken = redisToken
	}
	if err!=nil{
		return common.TokenDTO{},errors.New("redisToken非法"+err.Error())
	}
	webTokenObj,err := auth(webToken,userId)
	if err!=nil{
		return common.TokenDTO{},err
	}
	if webToken!=webTokenObj.Token{
		cResp := http.Cookie{
			Name: common.RedisToken,
			Value:webTokenObj.Token,
		}
		http.SetCookie(w,&cResp)
	}
	return webTokenObj,nil
}

//登录验证
func auth(webToken string,userId int64) (common.TokenDTO,error)  {
	if webToken =="" || strings.Replace(webToken," ","",-1)=="" || userId==0 {
		return common.TokenDTO{},errors.New("用户ID必传")
	}
	conn, err := common.ClientPool.Get()
	if err != nil {
		log.Fatal(err)
		return common.TokenDTO{},err
	}
	defer common.ClientPool.Put(conn)
	var replyToken common.TokenDTO
	//log.Printf(" webtoken %s \n " ,webToken)
	webTokenDTO,err:=common.AnalysisWebToken(webToken)
	if err != nil {
		log.Fatal(err)
		return common.TokenDTO{},err
	}
	timeScape :=webTokenDTO.LoginTime-time.Now().UnixMilli()
	var flagGetAndSet =false
	if timeScape>=common.TokenExpireTime30min{
		log.Fatal("登录时间超时，重新登录把")
		return common.TokenDTO{},errors.New("登录时间超时，重新登录把")
	}else if timeScape<common.TokenExpireTime30min && timeScape>=common.TokenExpireTime20min{
		webTokenDTO.LoginTime=webTokenDTO.LoginTime+common.TokenExpireTime20min
		flagGetAndSet=true
	}else if timeScape<common.TokenExpireTime20min && timeScape>=common.TokenExpireTime10min{
		webTokenDTO.LoginTime=webTokenDTO.LoginTime+common.TokenExpireTime10min
		flagGetAndSet=true
	}
	webTokenDTO.Token,err = common.BuildWebToken(webTokenDTO.UserId,webTokenDTO.UserName,webTokenDTO.LoginTime)
	if err != nil {
		log.Fatal(err)
		return common.TokenDTO{},err
	}
	redisTokenStr,err := common.BuildRedisToken(webTokenDTO.UserId,webTokenDTO.UserName)
	if err != nil {
		log.Fatal(err)
		return common.TokenDTO{},err
	}
	var redisTokenObj = CheckAuthDTO{Token: redisTokenStr,FlagGetAndSet:flagGetAndSet }
	redisTokenObjStr,_ :=Encode(redisTokenObj)
	err = conn.Call("Listener.RedisGetAndSetToken", redisTokenObjStr , &replyToken)
	if err!=nil{
		return common.TokenDTO{},err
	}
	if reflect.DeepEqual(replyToken,common.TokenDTO{}){
		return  common.TokenDTO{},errors.New("没有登录")
	}
	if userId!=replyToken.UserId {
		return common.TokenDTO{},errors.New("userid不对")
	}
	return webTokenDTO,nil
}
