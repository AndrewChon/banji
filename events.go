package banji

import "time"

const (
	Name = "banji"

	StartTopic    = Name + "Start"
	StopTopic     = Name + "Stop"
	PreTickTopic  = Name + "PreTick"
	PostTickTopic = Name + "PostTick"
	ErrorTopic    = Name + "Error"
)

type StartEvent struct {
	Event
}

func (s *StartEvent) Topic() string {
	return StartTopic
}

func newStartEvent() *StartEvent {
	return new(StartEvent)
}

type StopEvent struct {
	Event
}

func (s *StopEvent) Topic() string {
	return StopTopic
}

func newStopEvent() *StopEvent {
	return new(StopEvent)
}

type PreTickEvent struct {
	Event
	tick time.Time
}

func (t *PreTickEvent) Topic() string {
	return PreTickTopic
}

func (t *PreTickEvent) Tick() time.Time {
	return t.tick
}

func newPreTickEvent(tick time.Time) *PreTickEvent {
	return &PreTickEvent{tick: tick}
}

type PostTickEvent struct {
	Event
	tick time.Time
	end  time.Time
}

func (t *PostTickEvent) Topic() string {
	return PostTickTopic
}

func (t *PostTickEvent) Tick() time.Time {
	return t.tick
}

func (t *PostTickEvent) End() time.Time {
	return t.end
}

func newPostTickEvent(tick, end time.Time) *PostTickEvent {
	return &PostTickEvent{
		tick: tick,
		end:  end,
	}
}

type ErrorEvent struct {
	Event
	err error
}

func (e *ErrorEvent) Topic() string {
	return ErrorTopic
}

func (e *ErrorEvent) Error() error {
	return e.err
}

func NewErrorEvent(err error) *ErrorEvent {
	return &ErrorEvent{
		err: err,
	}
}
