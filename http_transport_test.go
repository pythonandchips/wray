package wray

import (
	"encoding/json"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsUsable(t *testing.T) {
	Convey("usable if the url is http", t, func() {
		var url string
		var httpTransport HttpTransport
		var isUsable = true
		Given(func() { url = "http://localhost" })
		Given(func() { httpTransport = HttpTransport{} })
		When(func() { isUsable = httpTransport.isUsable(url) })
		Then(func() { So(isUsable, ShouldEqual, true) })
	})
	Convey("usable if the url is https", t, func() {
		var url string
		var httpTransport HttpTransport
		var isUsable = true
		Given(func() { url = "https://localhost" })
		Given(func() { httpTransport = HttpTransport{} })
		When(func() { isUsable = httpTransport.isUsable(url) })
		Then(func() { So(isUsable, ShouldEqual, true) })
	})
	Convey("usable if the url is ws", t, func() {
		var url string
		var httpTransport HttpTransport
		var isUsable = true
		Given(func() { url = "ws://localhost" })
		Given(func() { httpTransport = HttpTransport{} })
		When(func() { isUsable = httpTransport.isUsable(url) })
		Then(func() { So(isUsable, ShouldEqual, false) })
	})
	Convey("usable if the url is not blank", t, func() {
		var url string
		var httpTransport HttpTransport
		var isUsable = true
		Given(func() { url = "" })
		Given(func() { httpTransport = HttpTransport{} })
		When(func() { isUsable = httpTransport.isUsable(url) })
		Then(func() { So(isUsable, ShouldEqual, false) })
	})
}

func TestConnectionType(t *testing.T) {
	Convey("connection type is long-polling", t, func() {
		var httpTransport HttpTransport
		Given(func() { httpTransport = HttpTransport{} })
		Then(func() { So(httpTransport.connectionType(), ShouldEqual, "long-polling") })
	})
}

func TestSendToServer(t *testing.T) {

	Convey("should send request to server and return response to callback", t, func() {
		var server *httptest.Server
		var httpTransport HttpTransport
		var params map[string]interface{}
		var recievedParams map[string]interface{}
		var response Response
		var headerParams map[string]interface{}
		var messageParams map[string]interface{}
		var returnData []interface{}
		Given(func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestData, _ := ioutil.ReadAll(r.Body)
				json.Unmarshal(requestData, &recievedParams)

				returnBytes, _ := json.Marshal(returnData)
				w.Write(returnBytes)
			}))
		})
		Given(func() {
			headerParams = map[string]interface{}{"id": "1",
				"clientId":   "client1",
				"channel":    "/meta/connect",
				"successful": true,
				"advice": map[string]interface{}{"reconnect": "retry",
					"interval": 0,
					"timeout":  45000}}
		})
		Given(func() {
			messageParams = map[string]interface{}{"channel": "/foo", "data": map[string]interface{}{"hello": "world"}, "id": "3"}
		})
		Given(func() { returnData = []interface{}{headerParams, messageParams} })
		Given(func() { httpTransport = HttpTransport{url: server.URL} })
		Given(func() { params = map[string]interface{}{"hello": "world"} })
		When(func() { response, _ = httpTransport.send(params) })
		Convey("received paramters", func() {
			Then(func() { So(recievedParams, ShouldResemble, params) })
		})
		Convey("return response", func() {
			Then(func() { So(response.id, ShouldEqual, "1") })
			Then(func() { So(response.clientId, ShouldEqual, "client1") })
			Then(func() { So(response.successful, ShouldEqual, true) })
			Then(func() { So(len(response.messages), ShouldEqual, 1) })
		})
		Convey("returned message", func() {
			var message Message
			Given(func() { message = response.messages[0] })
			Then(func() { So(message.Channel, ShouldEqual, "/foo") })
			Then(func() { So(message.Data, ShouldResemble, map[string]interface{}{"hello": "world"}) })
			Then(func() { So(message.Id, ShouldEqual, "3") })
		})
		defer server.Close()
	})
	Convey("should send request to server and return response to callback", t, func() {
		var server *httptest.Server
		var httpTransport HttpTransport
		var params map[string]interface{}
		var recievedParams map[string]interface{}
		var response Response
		var headerParams map[string]interface{}
		var returnData []interface{}
		Given(func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestData, _ := ioutil.ReadAll(r.Body)
				json.Unmarshal(requestData, &recievedParams)

				returnBytes, _ := json.Marshal(returnData)
				w.Write(returnBytes)
			}))
		})
		Given(func() {
			headerParams = map[string]interface{}{"id": "1",
				"clientId":   "client1",
				"channel":    "/meta/connect",
				"successful": true,
				"advice": map[string]interface{}{"reconnect": "retry",
					"interval": 0,
					"timeout":  45000}}
		})
		Given(func() { returnData = []interface{}{headerParams} })
		Given(func() { httpTransport = HttpTransport{url: server.URL} })
		Given(func() { params = map[string]interface{}{"hello": "world"} })
		When(func() { response, _ = httpTransport.send(params) })
		Then(func() { So(len(response.messages), ShouldEqual, 0) })
		defer server.Close()
	})
	Convey("when server is not available", t, func() {
		var httpTransport HttpTransport
		var params map[string]interface{}
		var response Response
		Given(func() { httpTransport = HttpTransport{url: "127.0.0.40"} })
		Given(func() { params = map[string]interface{}{"hello": "world"} })
		When(func() { response, _ = httpTransport.send(params) })
		Then(func() { So(response.successful, ShouldEqual, false) })
	})
	Convey("when the server does not return an OK response", t, func() {
		var server *httptest.Server
		var httpTransport HttpTransport
		var params map[string]interface{}
		var response Response
		var err error
		Given(func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.NotFound(w, r)
			}))
		})
		Given(func() { httpTransport = HttpTransport{url: server.URL} })
		Given(func() { params = map[string]interface{}{"hello": "world"} })
		When(func() { response, err = httpTransport.send(params) })
		Then(func() { So(err, ShouldResemble, errors.New("404 Not Found")) })
		Then(func() { So(response.successful, ShouldEqual, false) })
		defer server.Close()
	})
}
