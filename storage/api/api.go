package api

import (
	"context"
	"fmt"

	"github.com/Nestik55/task0/model"
	"github.com/Nestik55/task0/storage/cachik"
	"github.com/Nestik55/task0/storage/postgres"
)

var Cash = cachik.NewCache()
var DB = postgres.NewPostgres()

func GetOrder(key string) (model.Order, bool) {
	order, err := Cash.Get(key)
	if err != nil {
		fmt.Println(err)
		err = DB.GetOrder(context.TODO(), &order, key)
		if err != nil {
			fmt.Println(err)
			return model.Order{}, false
		}
	}

	return order, true
}

func SetOrder(key string, order model.Order) {
	//	fmt.Println("cash.set <- begin")
	Cash.Set(key, order)
	//fmt.Println("cash.set <- end")
	DB.SetOrder(context.TODO(), &order)
}
