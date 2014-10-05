package wray

import "time"

type Schedular interface {
	wait(time.Duration, func())
	delay() time.Duration
	sleep(time.Duration)
}

type ChannelSchedular struct {
}

func (self ChannelSchedular) sleep(delay time.Duration) {
	time.Sleep(delay)
}

func (self ChannelSchedular) wait(delay time.Duration, callback func()) {
	go func() {
		time.Sleep(delay)
		callback()
	}()
}

func (self ChannelSchedular) delay() time.Duration {
	return (1 * time.Minute)
}
