package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net/url"
	"os"
	"time"
)

type config struct {
	broker    string
	clientCnt int
}

func main() {
	c := &config{}
	fs := flag.NewFlagSet("mqtt-benchmark", flag.ExitOnError)
	fs.StringVar(&c.broker, "b", "", "Broker address, eg: localhost:1883.")
	fs.IntVar(&c.clientCnt, "c", 1, "Client count.")
	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}
	if len(c.broker) == 0 {
		panic("missing broker address")
	}

	resCh := make(chan string)
	for i := 0; i < c.clientCnt; i++ {
		time.Sleep(500 * time.Millisecond)
		go doClientConn(i, c.broker, resCh)
	}

	for i := 0; i < c.clientCnt; i++ {
		res := <-resCh
		fmt.Println(res)
	}
}

func doClientConn(num int, broker string, resCh chan string) {
	c := mqtt.NewClient(&mqtt.ClientOptions{
		Servers:        []*url.URL{{Host: broker, Scheme: "tcp"}},
		ClientID:       fmt.Sprintf("client@%d", num),
		Username:       "100014",
		Password:       "",
		KeepAlive:      60,
		PingTimeout:    10 * time.Second,
		ConnectTimeout: 10 * time.Second,
		DefaultPublishHandler: func(client mqtt.Client, message mqtt.Message) {
			if num == 1 {
				fmt.Println(string(message.Payload()))
			}
		},
		OnConnect: func(client mqtt.Client) {

		},
		OnConnectionLost: func(client mqtt.Client, err error) {

		},
		OnReconnecting: func(client mqtt.Client, options *mqtt.ClientOptions) {

		},
		OnConnectAttempt: func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {

			return tlsCfg
		},
		WriteTimeout: 10 * time.Second,
	})

	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
