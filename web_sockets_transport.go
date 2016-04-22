package wray

type WebSocketTransport struct {
}

func (self WebSocketTransport) isUsable(string) bool {
	return false
}
func (self WebSocketTransport) connectionType() string {
	return ""
}
func (self WebSocketTransport) send(map[string]interface{}) (Response, error) {
	return Response{}, nil
}
func (self WebSocketTransport) setUrl(string) {
}
