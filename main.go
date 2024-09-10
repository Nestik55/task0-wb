package main

import (
	"encoding/json"
	"fmt"

	"github.com/Nestik55/task0/model"
	"github.com/Nestik55/task0/publisher"
	"github.com/Nestik55/task0/server"
	"github.com/Nestik55/task0/storage/api"
	"github.com/nats-io/stan.go"
)

const (
	clusterID   = "test-cluster"
	clientID    = "client"
	channelName = "orderChannel"
)

func msgHandler(msg *stan.Msg) {
	order := model.Order{}
	if err := json.Unmarshal(msg.Data, &order); err != nil {
		fmt.Println(err)
	}
	api.SetOrder(order.OrderUID, order)
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

	publisher.Publish(stanConn, channelName, "publisher/orders/model1.json")
	//time.Sleep(2 * time.Second)
	server.Run()

	err = sub.Close()

	if err != nil {
		fmt.Println(err)
	}

}
