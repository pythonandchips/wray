package wray

import (
	. "github.com/smartystreets/goconvey/convey"
  "testing"
)

func TestSelectTransport(t *testing.T) {
  Convey("when all transports are usable", t, func(){
    var transportTypes []string
    var transport Transport
    var fakeHttpTransport *FakeHttpTransport
    var fayeClient *FayeClient

    Given(func(){ transportTypes = []string{"long-polling"} })
    Given(func(){ fakeHttpTransport = &FakeHttpTransport{usable: true} })
    Given(func(){ fayeClient = &FayeClient{url: "http://localhost"} })
    Given(func(){ registeredTransports = []Transport{fakeHttpTransport}   })
    When(func(){ transport, _ = SelectTransport(fayeClient, transportTypes, []string{}) })
    Then(func(){ So(transport, ShouldEqual, fakeHttpTransport) })
  })

  Convey("when no transports are usable", t, func(){
    var transportTypes []string
    var fakeHttpTransport *FakeHttpTransport
    var fayeClient *FayeClient
    var err error

    Given(func(){ transportTypes = []string{"long-polling"} })
    Given(func(){ fakeHttpTransport = &FakeHttpTransport{usable: false} })
    Given(func(){ fayeClient = &FayeClient{url: "http://localhost"} })
    Given(func(){ registeredTransports = []Transport{fakeHttpTransport}   })
    When(func(){ _, err = SelectTransport(fayeClient, transportTypes, []string{}) })
    Then(func(){ So(err, ShouldNotBeNil) })
  })
}
