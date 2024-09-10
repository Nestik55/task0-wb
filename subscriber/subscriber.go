package subscriber

import (
	"encoding/json"
	"fmt"

	"github.com/Nestik55/task0/model"
	"github.com/Nestik55/task0/storage/api"
	"github.com/nats-io/stan.go"
)

func msgHandler(msg *stan.Msg) {
	order := model.Order{}
	if err := json.Unmarshal(msg.Data, &order); err != nil {
		fmt.Println(err)
	}
	api.SetOrder(order.OrderUID, order)
	fmt.Println(order.OrderUID, order)
}
