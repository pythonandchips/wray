package wray

type FakeHttpTransport struct {
	usable     bool
	timesSent  int
	sentParams map[string]interface{}
	response   Response
	url        string
	err        error
}

func (self FakeHttpTransport) isUsable(endpoint string) bool {
	return self.usable
}

func (self FakeHttpTransport) connectionType() string {
	return "long-polling"
}

func (self *FakeHttpTransport) send(params map[string]interface{}) (Response, error) {
	self.sentParams = params
	self.timesSent++
	return self.response, self.err
}

func (self *FakeHttpTransport) setUrl(url string) {
	self.url = url
}
