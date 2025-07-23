# Banji

A concurrent event engine (more specifically, proactor) and EDA framework for Go.

>**Disclaimer:**
> This project is _very, very much_ in its early stages. This codebase is subject to significant change. Banji should be 
> considered unfinished at this time. Use at your own risk!

## About

Banji aims to be a simple way to implement a publish-subscribe pattern in your projects. Banji was written primarily for
my own use in some other projects of mine, but I figured it may serve a purpose for others as well.

## Features

- Built on a pairing heap priority queue
- Encourages decoupling of areas-of-concern
- Concurrent demultiplexing of events
- Graceful shutdown mechanism
- Tick-based engine loop
- Reflection-free routing
- Built-in events

## How To Use

Implementing the engine is relatively straight-forward:

```go
func main() {
    eng := banji.New(
        banji.WithTPS(64),
        banji.WithDemuxers(8),
    )

// Add your receivers ...

    eng.Start()
    defer eng.Stop()
}
```

Events and Receivers are defined according to this pattern:

```go
const MyTopic string = "mycomponent.My"

type MyEvent struct {
    banji.Event
    // ... Data ...
}

func (e *MyEvent) Topic() string {
    return MyTopic
}

type MyReceiver struct {
    banji.Receiver
    // ... Data ...
}

func (r *MyReceiver) Topic() string {
    return MyTopic
}

func (r *MyReceiver) Handle(em banji.Emitter) error {
    e, _ := em.(*MyEvent)
    // ... Handle the event
    return nil
}
```