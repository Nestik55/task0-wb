package publisher

import (
	"fmt"
	"io"
	"os"

	"github.com/nats-io/stan.go"
)

func Publish(stanConn stan.Conn, channelName, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	err = stanConn.Publish(channelName, data)
	if err != nil {
		fmt.Println(err)
	}
}
