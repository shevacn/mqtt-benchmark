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
	broker           string
	clientCnt        int
	printReceivedMsg bool
	sleepMs          int
}

func main() {
	c := &config{}
	fs := flag.NewFlagSet("mqtt-benchmark", flag.ExitOnError)
	fs.StringVar(&c.broker, "b", "", "Broker address, eg: localhost:1883.")
	fs.IntVar(&c.clientCnt, "c", 1, "Client count.")
	fs.IntVar(&c.sleepMs, "s", 1, "Sleep ms.")
	fs.BoolVar(&c.printReceivedMsg, "p", false, "Print received msg.")
	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}
	if len(c.broker) == 0 {
		panic("missing broker address")
	}

	resCh := make(chan string)
	for i := 0; i < c.clientCnt; i++ {
		time.Sleep(time.Duration(c.sleepMs) * time.Millisecond)
		go doClientConn(i, c, resCh)
	}

	for i := 0; i < c.clientCnt; i++ {
		res := <-resCh
		fmt.Println(res)
	}
}

func doClientConn(num int, conf *config, resCh chan string) {
	c := mqtt.NewClient(&mqtt.ClientOptions{
		Servers:        []*url.URL{{Host: conf.broker, Scheme: "tcp"}},
		ClientID:       fmt.Sprintf("client@%d", num),
		Username:       "100014",
		Password:       "",
		KeepAlive:      60,
		PingTimeout:    10 * time.Second,
		ConnectTimeout: 10 * time.Second,
		DefaultPublishHandler: func(client mqtt.Client, message mqtt.Message) {
			if conf.printReceivedMsg && num == 0 {
				fmt.Println(string(message.Payload()))
			}
		},
		OnConnect: func(client mqtt.Client) {
		},
		OnConnectionLost: func(client mqtt.Client, err error) {
			fmt.Printf("[client#%d]Connection lost\n", num)
		},
		OnReconnecting: func(client mqtt.Client, options *mqtt.ClientOptions) {
			fmt.Printf("[client#%d]Reconnecting\n", num)

		},
		OnConnectAttempt: func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
			fmt.Printf("[client#%d]Connect attempt\n", num)

			return tlsCfg
		},
		WriteTimeout: 10 * time.Second,
	})

	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
