package common

import (
	pool "github.com/flyaways/pool"
	"log"
	"time"
)

var ClientPool *pool.RPCPool

func ClientPoolInit()  {
	options := &pool.Options{
		InitTargets:  []string{"127.0.0.1:12345"},
		InitCap:      30,
		MaxCap:       30000,
		DialTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 60,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
	}
	var err error
	ClientPool, err = pool.NewRPCPool(options) 			//for rpc
	//p, err := pool.NewTCPPool(options)			//for tcp
	if err != nil {
		log.Printf("%#v\n", err)
		return
	}
	if ClientPool == nil {
		log.Printf("p= %#v\n", ClientPool)
		return
	}
}