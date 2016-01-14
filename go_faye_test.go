package wray

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func Given(execute func()) {
	execute()
}
func When(execute func()) {
	execute()
}
func Then(execute func()) {
	execute()
}

func Wait(delay time.Duration) {
	time.Sleep(delay)
}

func Pending(description string, execute func()) {
	fmt.Println("PENDING: " + description)
}

func TestInitializingClient(t *testing.T) {
	Convey("puts the client in the unconnected state", t, func() {
		var fayeClient *FayeClient
		When(func() { fayeClient = NewFayeClient("http://localhost") })
		Then(func() { So(fayeClient.state, ShouldEqual, UNCONNECTED) })
		Then(func() { So(fayeClient.schedular, ShouldHaveSameTypeAs, ChannelSchedular{}) })
	})
}

func TestSubscribe(t *testing.T) {
	Convey("subscribe to a channel when unconnected", t, func() {
		var fayeClient FayeClient
		var callback func(Message)
		var subscriptionPromise SubscriptionPromise
		var fakeHttpTransport *FakeHttpTransport
		var subscriptionParams map[string]interface{}
		var response Response
		var err error
		Given(func() {
			response = Response{id: "1", channel: "/meta/handshake", successful: true, clientId: "client4", supportedConnectionTypes: []string{"long-polling"}}
		})
		Given(func() { fakeHttpTransport = &FakeHttpTransport{usable: true, response: response} })
		Given(func() { registeredTransports = []Transport{fakeHttpTransport} })
		Given(func() { fayeClient = BuildFayeClient().WithTransport(fakeHttpTransport).Client() })
		Given(func() {
			subscriptionParams = map[string]interface{}{"channel": "/meta/subscribe", "clientId": response.clientId, "subscription": "/foo/*", "id": "1"}
		})
		Given(func() { callback = func(message Message) {} })
		When(func() { subscriptionPromise, err = fayeClient.Subscribe("/foo/*", false, callback) })
		Convey("connects the faye client", func() {
			Then(func() { So(err, ShouldEqual, nil) })
			Then(func() { So(fayeClient.state, ShouldEqual, CONNECTED) })
		})
		Convey("add the subscription to the client", func() {
			Then(func() { So(len(fayeClient.subscriptions), ShouldEqual, 1) })
			Then(func() { So(fayeClient.subscriptions[0].channel, ShouldEqual, subscriptionPromise.subscription.channel) })
			Then(func() {
				So(fayeClient.subscriptions[0].callback, ShouldEqual, subscriptionPromise.subscription.callback)
			})
		})
		Convey("the promise has the setup subscription", func() {
			Then(func() { So(subscriptionPromise.subscription.channel, ShouldEqual, "/foo/*") })
			Then(func() { So(subscriptionPromise.subscription.callback, ShouldEqual, callback) })
		})
		Convey("the client send the subscription to the server", func() {
			Then(func() { So(fakeHttpTransport.sentParams, ShouldResemble, subscriptionParams) })
		})
	})
}

func TestSubscriptionError(t *testing.T) {
	Convey("subscribe to a channel when unconnected", t, func() {
		var fayeClient FayeClient
		var callback func(Message)
		var subscriptionPromise SubscriptionPromise
		var fakeHttpTransport *FakeHttpTransport
		var subscriptionParams map[string]interface{}
		var failedResponse Response
		var clientId string = "client1"
		var err error
		Given(func() {
			failedResponse = Response{id: "1", channel: "/meta/subscribe", successful: false, clientId: clientId, supportedConnectionTypes: []string{"long-polling"}}
		})
		Given(func() { fakeHttpTransport = &FakeHttpTransport{usable: true, response: failedResponse} })
		Given(func() { registeredTransports = []Transport{fakeHttpTransport} })
		Given(func() { fayeClient = BuildFayeClient().WithTransport(fakeHttpTransport).Client() })
		Given(func() { fayeClient.state = CONNECTED })
		Given(func() {
			subscriptionParams = map[string]interface{}{"channel": "/meta/subscribe", "clientId": clientId, "subscription": "/foo/*", "id": "1"}
		})
		Given(func() { callback = func(message Message) {} })
		When(func() { subscriptionPromise, err = fayeClient.Subscribe("/foo/*", false, callback) })
		Convey("fails to subscribe", func() {
			Then(func() { So(err, ShouldNotEqual, nil) })
			Then(func() { So(subscriptionPromise.Successful(), ShouldEqual, false) })
		})
		Convey("not add the subscription to the client", func() {
			Then(func() { So(len(fayeClient.subscriptions), ShouldEqual, 0) })
		})
		Convey("the client send the subscription to the server", func() {
			Then(func() { So(fakeHttpTransport.sentParams, ShouldResemble, subscriptionParams) })
		})
	})
}

func TestPerformHandshake(t *testing.T) {
	Convey("successful handshake with server", t, func() {
		var fayeClient FayeClient
		var fakeHttpTransport *FakeHttpTransport
		var handshakeParams map[string]interface{}
		var response Response
		Given(func() {
			handshakeParams = map[string]interface{}{"channel": "/meta/handshake",
				"version":                  "1.0",
				"supportedConnectionTypes": []string{"long-polling"}}
		})

		Given(func() {
			response = Response{id: "1", channel: "/meta/handshake", successful: true, clientId: "client4", supportedConnectionTypes: []string{"long-polling"}}
		})
		Given(func() { fakeHttpTransport = &FakeHttpTransport{usable: true, response: response} })
		Given(func() { registeredTransports = []Transport{fakeHttpTransport} })
		Given(func() { fayeClient = BuildFayeClient().Client() })
		When(func() { fayeClient.handshake() })
		Then(func() { So(fayeClient.state, ShouldEqual, CONNECTED) })
		Then(func() { So(fayeClient.transport, ShouldEqual, fakeHttpTransport) })
		Then(func() { So(fayeClient.clientId, ShouldEqual, "client4") })
		Then(func() { So(fakeHttpTransport.sentParams, ShouldResemble, handshakeParams) })
		Then(func() { So(fakeHttpTransport.url, ShouldEqual, fayeClient.url) })
	})

	Convey("unsuccessful handshake with server", t, func() {
		var fayeClient FayeClient
		var fakeHttpTransport *FakeHttpTransport
		var handshakeParams map[string]interface{}
		var response Response
		Given(func() {
			handshakeParams = map[string]interface{}{"channel": "/meta/handshake",
				"version":                  "1.0",
				"supportedConnectionTypes": []string{"long-polling"}}
		})

		Given(func() {
			response = Response{id: "1", channel: "/meta/handshake", successful: false, clientId: "client4", supportedConnectionTypes: []string{"long-polling"}}
		})
		Given(func() {
			fakeHttpTransport = &FakeHttpTransport{usable: true, response: response, err: errors.New("it didny work")}
		})
		Given(func() { registeredTransports = []Transport{fakeHttpTransport} })
		Given(func() { fayeClient = BuildFayeClient().Client() })
		When(func() { fayeClient.handshake() })
		Then(func() { So(fayeClient.state, ShouldEqual, UNCONNECTED) })
		Then(func() { So(fayeClient.schedular.delay(), ShouldEqual, 10*time.Second) })
	})

	Convey("handshake with no available transports", t, func() {
		var fayeClient FayeClient
		var fakeHttpTransport *FakeHttpTransport
		Given(func() { fakeHttpTransport = &FakeHttpTransport{usable: false} })
		Given(func() { registeredTransports = []Transport{fakeHttpTransport} })
		Given(func() { fayeClient = BuildFayeClient().Client() })
		Then(func() { So(func() { fayeClient.handshake() }, ShouldPanicWith, "No usable transports available") })
	})
	Convey("when server does not support available transports", t, func() {
		var fayeClient FayeClient
		var fakeHttpTransport *FakeHttpTransport
		var handshakeParams map[string]interface{}
		var response Response
		Given(func() {
			handshakeParams = map[string]interface{}{"channel": "/meta/handshake",
				"version":                  "1.0",
				"supportedConnectionTypes": []string{"long-polling"}}
		})

		Given(func() {
			response = Response{id: "1", channel: "/meta/handshake", successful: true, clientId: "client4", supportedConnectionTypes: []string{"web-socket"}}
		})
		Given(func() { fakeHttpTransport = &FakeHttpTransport{usable: true, response: response} })
		Given(func() { registeredTransports = []Transport{fakeHttpTransport} })
		Given(func() { fayeClient = BuildFayeClient().Client() })
		Then(func() {
			So(func() { fayeClient.handshake() }, ShouldPanicWith, "Server does not support any available transports. Supported transports: web-socket")
		})
	})
}

func TestHandleResponse(t *testing.T) {
	Convey("when server does not support available transports", t, func() {
		var fayeClient FayeClient
		var response Response
		var messages []Message
		var subscriptions []Subscription
		var firstParams map[string]interface{}
		var secondParams map[string]interface{}
		var firstMessages []map[string]interface{}
		var secondMessages []map[string]interface{}
		Given(func() {
			subscriptions = []Subscription{
				{"/foo/bar", func(message Message) { firstMessages = append(firstMessages, message.Data) }},
				{"/foo/*", func(message Message) { secondMessages = append(secondMessages, message.Data) }},
			}
		})
		Given(func() { firstParams = map[string]interface{}{"foo": "bar"} })
		Given(func() { secondParams = map[string]interface{}{"baz": "qux"} })
		Given(func() { fayeClient = BuildFayeClient().WithSubscriptions(subscriptions).Client() })
		Given(func() {
			messages = []Message{
				{Channel: "/foo/bar", Id: "1", Data: firstParams},
				{Channel: "/foo/quz", Id: "1", Data: secondParams},
			}
		})
		Given(func() { response = Response{messages: messages, channel: "/meta/connect", clientId: "client1"} })
		When(func() { fayeClient.handleResponse(response) })
		//need a very short sleep in here to allow the go routines to complete
		//as all they are doing is assigning a variable 10 milliseconds shoule be more than enough
		Wait(100 * time.Millisecond)
		Then(func() { So(firstMessages, ShouldContain, firstParams) })
		Then(func() { So(len(secondMessages), ShouldEqual, 2) })
		Then(func() { So(secondMessages, ShouldContain, firstParams) })
		Then(func() { So(secondMessages, ShouldContain, secondParams) })
	})
}

func TestPublish(t *testing.T) {
	Convey("publish message to server", t, func() {
		var fayeClient FayeClient
		var fakeHttpTransport *FakeHttpTransport
		var data map[string]interface{}
		var response Response
		Given(func() {
			response = Response{id: "1", channel: "/meta/handshake", successful: true, clientId: "client4", supportedConnectionTypes: []string{"long-polling"}}
		})
		Given(func() { fakeHttpTransport = &FakeHttpTransport{usable: true, response: response} })
		Given(func() { registeredTransports = []Transport{fakeHttpTransport} })
		Given(func() { fayeClient = BuildFayeClient().WithTransport(fakeHttpTransport).Connected().Client() })
		Given(func() { data = map[string]interface{}{"hello": "world"} })
		When(func() { fayeClient.Publish("/foo", data) })
		Then(func() {
			So(fakeHttpTransport.sentParams, ShouldResemble, map[string]interface{}{"channel": "/foo", "data": data, "clientId": fayeClient.clientId})
		})

	})
}
