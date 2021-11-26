package test

import (
	"fmt"
	"gotest1/src/common"
	"testing"
)

func Test1(t *testing.T)  {
	//加密数据
	result,_:=common.DesEncrypt_CBC([]byte("zhangbao_123123123"))
	fmt.Println(result)

	//解密
	result,_=common.DesDecrypt_CBC(result)
	fmt.Println("解密之后的数据:",string(result))
}
