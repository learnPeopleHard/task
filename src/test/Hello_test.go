package test

import (
	"fmt"
	"testing"
)

//声明全局变量
var c int  = 100

func TestHello(t *testing.T)  {
	var a int = 100
	var b int = 200
	b, a = a, b
	fmt.Println(a, b,c)

	fmt.Println("hello world")
}