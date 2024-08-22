package main

import (
	"fmt"
	"strconv"

	"github.com/Nestik55/task0/cache"
	"github.com/nats-io/stan.go"
)

const (
	clusterID   = "test-cluster"
	clientID    = "client"
	channelName = "orderChannel"
)

var cash = cache.NewCache()

var j = 0

func msgHandler(msg *stan.Msg) {
	cash.Set(strconv.Itoa(j), string(msg.Data))
	fmt.Println(string(msg.Data))
	j++
}

func main() {

	stanConn, err := stan.Connect(clusterID, clientID)
	if err != nil {
		panic(err)
	}

	sub, err := stanConn.Subscribe(channelName, msgHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("123")
	stanConn.Publish(channelName, []byte("Hello"))
	stanConn.Publish(channelName, []byte("Hello2"))
	stanConn.Publish(channelName, []byte("Hello3"))
	fmt.Println("456")
	fmt.Println(cash.Get(strconv.Itoa(j - 1)))
	fmt.Println(cash.Get(strconv.Itoa(j - 2)))
	fmt.Println(cash.Get(strconv.Itoa(j - 3)))

	err = sub.Close()
	if err != nil {
		fmt.Println(err)
	}

}
