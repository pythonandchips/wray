package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
)

func main() {
	origin := "http://localhost:5000/"
	url := "ws://localhost:5000/faye"
	fmt.Println("Dialing....")
	ws, err := websocket.Dial(url, "", origin)
	fmt.Println("Dialled")
	if err != nil {
		panic(err)
	}
	handshakeParams := map[string]interface{}{"channel": "/meta/handshake",
		"version":                  "1.0",
		"supportedConnectionTypes": []string{"websocket"}}

	writeToWs(ws, handshakeParams)

	msg, n := readFromWs(ws)
	var jsonData []interface{}
	json.Unmarshal(msg[:n], &jsonData)
	fmt.Println(jsonData)
	headerData := jsonData[0].(map[string]interface{})
	var clientId string
	if headerData["clientId"] != nil {
		clientId = headerData["clientId"].(string)
	}
	subscriptionParams := map[string]interface{}{"channel": "/meta/subscribe", "clientId": clientId, "subscription": "/foo", "id": "1"}
	writeToWs(ws, subscriptionParams)
	readFromWs(ws)
	connectParams := map[string]interface{}{"channel": "/meta/connect", "clientId": clientId, "connectionType": "web-socket"}
	writeToWs(ws, connectParams)
	readFromWs(ws)
	for {
		readFromWs(ws)
	}
}

func readFromWs(ws *websocket.Conn) ([]byte, int) {
	fmt.Println("Reading")
	var r int
	var subread error
	var readBytes = make([]byte, 2048)
	if r, subread = ws.Read(readBytes); subread != nil {
		fmt.Println("panic on reading")
		panic(subread)
	}
	fmt.Printf("Read: %s.\n", readBytes[:r])
	return readBytes, r
}

func writeToWs(ws *websocket.Conn, message map[string]interface{}) {
	dataBytes, _ := json.Marshal(message)
	fmt.Println("Writing....")
	if _, err := ws.Write(dataBytes); err != nil {
		panic(err)
	}
	fmt.Println("Written")
}
