package wray

import "time"

type FakeSchedular struct {
	requestedDelay time.Duration
}

func (self *FakeSchedular) sleep(delay time.Duration) {
	self.requestedDelay = delay
}

func (self *FakeSchedular) wait(delay time.Duration, callback func()) {
	self.requestedDelay = delay
}

func (self *FakeSchedular) delay() time.Duration {
	return self.requestedDelay
}

type FayeClientBuilder struct {
	client FayeClient
}

func BuildFayeClient() FayeClientBuilder {
	schedular := &FakeSchedular{}
	client := FayeClient{state: UNCONNECTED, url: "https://localhost", clientId: "client1", schedular: schedular}
	return FayeClientBuilder{client}
}

func (self FayeClientBuilder) Client() FayeClient {
	return self.client
}

func (self FayeClientBuilder) Connected() FayeClientBuilder {
	self.client.state = CONNECTED
	return self
}

func (self FayeClientBuilder) WithTransport(transport Transport) FayeClientBuilder {
	self.client.transport = transport
	return self
}

func (self FayeClientBuilder) AddSubscription(subscription Subscription) FayeClientBuilder {
	self.client.subscriptions = append(self.client.subscriptions, subscription)
	return self
}

func (self FayeClientBuilder) WithSubscriptions(subscriptions []Subscription) FayeClientBuilder {
	self.client.subscriptions = subscriptions
	return self
}
