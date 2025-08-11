package test

import (
	"reflect"
	"testing"

	"github.com/AndrewChon/banji"

	"github.com/google/uuid"
)

/* banji.start */

type StartReceiverTest struct {
	banji.ReceiverEmbed
	t *testing.T
}

func (r *StartReceiverTest) Topic() string {
	return banji.StartTopic
}

func (r *StartReceiverTest) Handle(e banji.Event) error {
	event, ok := e.(*banji.StartEvent)
	if !ok {
		r.t.Fatalf("Expected *banji.StartEvent, received %v\n", reflect.TypeOf(event))
	}

	validateEventImplementation(event, r.t)
	return nil
}

/* banji.stop */

type StopReceiverTest struct {
	banji.ReceiverEmbed
	t *testing.T
}

func (r *StopReceiverTest) Topic() string {
	return banji.StopTopic
}

func (r *StopReceiverTest) Handle(e banji.Event) error {
	event, ok := e.(*banji.StopEvent)
	if !ok {
		r.t.Fatalf("Expected *banji.StopEvent, received %v\n", reflect.TypeOf(event))
	}

	validateEventImplementation(event, r.t)
	return nil
}

/* banji.preTick */

type PreTickReceiverTest struct {
	banji.ReceiverEmbed
	t *testing.T
}

func (r *PreTickReceiverTest) Topic() string {
	return banji.PreTickTopic
}

func (r *PreTickReceiverTest) Handle(e banji.Event) error {
	event, ok := e.(*banji.PreTickEvent)
	if !ok {
		r.t.Fatalf("Expected *banji.PreTickEvent, received %v\n", reflect.TypeOf(event))
	}

	validateEventImplementation(event, r.t)
	return nil
}

/* banji.postTick */

type PostTickReceiverTest struct {
	banji.ReceiverEmbed
	t *testing.T
}

func (r *PostTickReceiverTest) Topic() string {
	return banji.PostTickTopic
}

func (r *PostTickReceiverTest) Handle(e banji.Event) error {
	event, ok := e.(*banji.PostTickEvent)
	if !ok {
		r.t.Fatalf("Expected *banji.PostTickEvent, received %v\n", reflect.TypeOf(event))
	}

	validateEventImplementation(event, r.t)
	return nil
}

func validateEventImplementation(e banji.Event, t *testing.T) {
	if e.ID() == uuid.Nil {
		t.Fatalf("Posted Event %v has a nil ID\n", reflect.TypeOf(e))
	}

	if e.Topic() == "" {
		t.Fatalf("Event %v has an empty topic string\n", reflect.TypeOf(e))
	}

	if e.Postmark().IsZero() {
		t.Fatalf("Posted Event %v has an nil postmark\n", reflect.TypeOf(e))
	}
}
