package test

import (
	"fmt"
	"github.com/msterzhang/gpool"
	"gotest1/src/common"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func Test_test(b *testing.T)  {
	size:=100
	pool := gpool.New(size)
	for i :=0;i<200;i++{
		pool.Add(1)
		go runtiest(
			i,
			pool,
		)
	}
	pool.Wait()
}

func Test_sync(b *testing.T)  {
	params := url.Values{}
	Url, err := url.Parse("http://localhost:8081/login")
	if err != nil {
		return
	}
	params.Set("userName","zhangsan"+strconv.Itoa(rand.Intn(200)))
	params.Set("passWord","123456")
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath) // https://httpbin.org/get?age=23&name=zhaofan
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp,err := client.Get(urlPath)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func runtiest(i int,pool *gpool.Pool)  {
	for k :=0;k<10000000;k++{
		common.Try(func() {
			runTest(i,pool)
		}, func(err interface{}) {
			fmt.Printf("requst fail %v \n",err)
		})
	}
	pool.Done()
}

func runTest(i int,pool *gpool.Pool)  {

	params := url.Values{}
	Url, err := url.Parse("http://localhost:8081/login")
	if err != nil {
		return
	}
	params.Set("userName","zhangsan"+strconv.Itoa(i))
	params.Set("passWord","123456")
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Printf("url: %s \n",urlPath) // https://httpbin.org/get?age=23&name=zhaofan
	timeout := 5 * time.Second
	client := http.Client{
		Timeout: timeout,
	}
	resp,err := client.Get(urlPath)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("返回数据: %s \n", string(body))
}
