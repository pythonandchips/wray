package wray

type FakeHttpTransport struct {
  usable bool
  timesSent int
  sentParams map[string]interface{}
  response Response
  url string
}

func(self FakeHttpTransport) isUsable(endpoint string) bool {
  return self.usable
}

func(self FakeHttpTransport) connectionType() string {
  return "long-polling"
}

func(self *FakeHttpTransport) send(params map[string]interface{}, callback func(Response)) {
  self.sentParams = params
  callback(self.response)
  self.timesSent++
}

func(self *FakeHttpTransport) setUrl(url string) {
  self.url = url
}
