package common

import (
	"log"
	"strconv"
	"strings"
)

const TokenExpireTime10min = 10*60*1000
const TokenExpireTime20min = 20*60*1000
const TokenExpireTime30min = 30*60*1000


func BuildWebToken(UserId int64,name string,nowTime int64) (string,error)  {
	userStr := strconv.FormatInt(UserId,10)+"_"+strconv.FormatInt(nowTime,10)+"_"+name
	desUserStr,err := DesEncrypt_CBC([]byte(userStr))
	if err!=nil{
		log.Printf("UserId %s name %s 生产token失败 %s ",UserId, name,err.Error())
		return "",err
	}
	return desUserStr, nil
}

func BuildRedisToken(UserId int64,name string) (string,error)  {
	userStr := strconv.FormatInt(UserId,10)+"_"+name
	desUserStr,err := DesEncrypt_CBC([]byte(userStr))
	if err!=nil{
		log.Printf("UserId %s name %s 生产token失败 %s ",UserId, name,err.Error())
		return "",err
	}
	return desUserStr, nil
}


func AnalysisWebToken(token string) (TokenDTO,error)  {
	tokenStr,err := DesDecrypt_CBC(token)
	if err!=nil{
		return TokenDTO{},err
	}
	tokenObj := TokenDTO{}
	userIdStr :=strings.Split(tokenStr,"_")[0]
	tokenObj.UserId,err =strconv.ParseInt(userIdStr,10,64)
	if err!=nil{
		return TokenDTO{},err
	}
	tokenObj.Token = token
	expireStr := strings.Split(tokenStr,"_")[1]
	tokenObj.LoginTime,err=strconv.ParseInt(expireStr,10,64)
	if err!=nil{
		return TokenDTO{},err
	}
	tokenObj.UserName = tokenStr[len(userIdStr)+len(expireStr)+2:len(tokenStr)]
	return tokenObj, nil
}


