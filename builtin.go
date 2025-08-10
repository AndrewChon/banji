package banji

import (
	"time"
)

/* banji.start */

const StartTopic = "banji.start"

// StartEvent is an Event posted when the Engine starts. It primarily serves as a signal to components to perform any
// bootstrapping or initialization they may require.
type StartEvent struct {
	EventEmbed
}

func (e *StartEvent) Topic() string {
	return StartTopic
}

/* banji.stopLoop */

const StopTopic = "banji.stopLoop"

// StopEvent is an Event posted when the engine is about to shut down. It primarily serves as a signal to components to
// perform any shutdown or cleanup mechanics they may require. StopEvent should not be used as a catalyst for posting
// additional events, as they will not be routed by the engine once StopEvent has been posted.
type StopEvent struct {
	EventEmbed
}

func (e *StopEvent) Topic() string {
	return StopTopic
}

/* banji.preTick */

const PreTickTopic = "banji.preTick"

// PreTickEvent is an Event posted when a new tick has begun.
type PreTickEvent struct {
	EventEmbed
	tick time.Time
}

func (e *PreTickEvent) Topic() string {
	return PreTickTopic
}

func (e *PreTickEvent) Tick() time.Time {
	return e.tick
}

/* banji.postTick */

const PostTickTopic = "banji.postTick"

// PostTickEvent is an Event posted once the work for a given tick has concluded. Therefore, it does not necessarily
// mean that a new tick has also begun (you should listen for PreTickEvent instead for such purposes).
type PostTickEvent struct {
	EventEmbed
	tick time.Time
}

func (e *PostTickEvent) Topic() string {
	return PostTickTopic
}

func (e *PostTickEvent) Tick() time.Time {
	return e.tick
}

/* banji.error */

const ErrorTopic = "banji.error"

// ErrorEvent is an Event posted when a Receiver's handler method returns an error.
type ErrorEvent struct {
	EventEmbed
	err error
}

func (e *ErrorEvent) Topic() string {
	return ErrorTopic
}

func (e *ErrorEvent) Error() error {
	return e.err
}
