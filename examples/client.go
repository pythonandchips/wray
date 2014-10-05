package main

import "github.com/pythonandchips/wray"
import "fmt"

func main() {
	wray.RegisterTransports([]wray.Transport{&wray.HttpTransport{}})
	client := wray.NewFayeClient("http://localhost:5000/faye")

	fmt.Println("subscribing")
	client.Subscribe("/foo", false, func(message wray.Message) {
		fmt.Println("---receiving foo----------------------------------------")
		fmt.Println(message.Data)
	})
	client.Subscribe("/bar", false, func(message wray.Message) {
		fmt.Println("---receiving bar----------------------------------------")
		fmt.Println(message.Data)
	})

	client.Listen()
}
