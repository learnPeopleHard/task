package test

import (
	"container/list"
	"fmt"
	"strconv"
	"testing"
)
const pi = 3.14159 // 相当于 math.Pi 的近似值

type Week int

const (
	ONE Week =iota
	TWO
	THREE
	FOUR

)

func (c Week) String() string {
	switch c {
	case ONE:
		return "1"
	case TWO:
		return "2"
	default:
		return "N?A"
	}
}

func TestHello2(t *testing.T)  {
	var str =`123213
123123
123
123`
	var a int = Pi2
	var nn int = Pi22
	var b int = 200
	b, a = a, b
	fmt.Println(a, b,str,nn)

	var arr =[3]int {1,2,3}

	for i,v:=range arr {
		fmt.Printf("%d %d\n", i, v)
	}
	for _,v:=range arr {
		fmt.Printf("%d\n", v)
	}

	var arr1 =[3]int {}
	for i,v:=range arr1 {
		fmt.Printf("%d %d\n", i, v)
	}

	sence :=make(map[string]int)
	sence["a"] = 1
	sence["b"] = 2
	sence["c"] = 3

	delete(sence,"a")

	for k,v:=range sence{
		var str =k+"_"+strconv.Itoa(v)
		fmt.Println(str)
	}

	l:=list.New()
	l.PushBack("asd")
	l.PushBack("123")
	l.PushBack(43)
	for k := l.Front();k!=nil;k= k.Next(){
		fmt.Println(k.Value)
	}


	fmt.Println("hello world",arr[0])
}

type Cat struct {
	Color string
	Name  string
}
func NewCatByName(name string) *Cat {
	return &Cat{
		Name: name,
	}
}

func NewCatByColor(color string) *Cat {
	return &Cat{
		Color: color,
	}
}

type BlackCat struct {
	Cat  // 嵌入Cat, 类似于派生
}

// “构造子类”
func NewBlackCat(color string) *BlackCat {
	cat := &BlackCat{}
	cat.Color = color
	return cat
}
