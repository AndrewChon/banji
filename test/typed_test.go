package hello_world

import (
	"banji"
	"fmt"
	"testing"
	"time"
)

const GreetingTopic = "hello_world.Greeting"

type GreetingEvent struct {
	banji.Event
	target string
}

func (e *GreetingEvent) Topic() string {
	return GreetingTopic
}

func (e *GreetingEvent) Target() string {
	return e.target
}

type GreetingReceiver struct {
	banji.Receiver
	greeting string
}

func (r *GreetingReceiver) Topic() string {
	return GreetingTopic
}

func (r *GreetingReceiver) Greeting() string {
	return r.greeting
}

func (r *GreetingReceiver) Handle(em banji.Emittable) error {
	e, _ := em.(*GreetingEvent)
	fmt.Printf("%s, %s!\n", r.greeting, e.Target())
	return nil
}

func NewGreetingEvent(target string) *GreetingEvent {
	return &GreetingEvent{target: target}
}

func NewGreetingReceiver(greeting string) *GreetingReceiver {
	return &GreetingReceiver{greeting: greeting}
}

func TestHelloWorld(_ *testing.T) {
	eng := banji.New(
		banji.WithTPS(32),
		banji.WithDemuxers(8),
	)

	eng.AddReceiver(NewGreetingReceiver("Hello"))

	eng.Start()
	defer eng.Stop()

	eng.Post(NewGreetingEvent("World"), banji.Medium)

	time.Sleep(time.Second)
}
