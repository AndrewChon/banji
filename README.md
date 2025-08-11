# Banji

[![Go Report Card](https://goreportcard.com/badge/github.com/AndrewChon/banji)](https://goreportcard.com/report/github.com/AndrewChon/banji)
[![GoDoc](https://godoc.org/github.com/AndrewChon/banji?status.png)](https://godoc.org/github.com/AndrewChon/banji)

Banji is a concurrent event engine and EDA framework for Go.

> [!CAUTION]
> Banji is in the very, _very_ early stages of development.
> \
> There is undoubtedly room for improvement, and this codebase is bound to change. Please keep this in mind if you
> intend to use this package.

## About

I created Banji mainly for myself. Most packages that use an event engine, or are based on a pub-sub pattern, require
you to implement an interface with all of its listeners (e.g., gnet and gRPC). For most people, this is perfectly fine!
And these have become invaluable tools for many in the Go community. However, I wanted my projects to have a unified
event engine, from which many components can communicate with each other without needing to implement every possible
listener, allowing me to build onto it over time easily. In addition, I wanted a central, structured logging system for
all these components that would enable me to handle their logs in whichever way I wanted.

There are some events we do not need to handle or care about; some capabilities we simply do not need. Conversely, there
are some components we may desire in the future but cannot reasonably foresee needing now. This is the philosophy behind
Banji.

## Features

- Concurrent demultiplexing
- Event priorities
- Built-in events for happenings within the engine
- Proactor-esque design (see `ErrorEvent`)
- Reflection-free
- Graceful shutdown mechanics
- Tick-based engine loop
- Unified logging system

## Quick Start

Running the engine is pretty straight-forward.

```go
func main() {
    // This allows us to check for the program exiting.
    stop := make(chan os.Signal)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    eng := banji.New(
        banji.WithTPS(128),
        banji.WithDemuxers(8))
    )

    // Register your receivers here.
    eng.Subscribe(NewMyReceiver("Hiii!"))

    eng.Start() // This is non-blocking.
    defer eng.Stop()

    // Once the engine is running, you and other components
    // can then post events to the engine.
    eng.Post(NewMyEvent("Hello!"), 255)
    
    <- stop
}
```

Events implement the `banji.Event` interface by embedding `banji.EventEmbed`.

```go
// You should define a constant for your topic strings, for everyone else's sake.

const MyTopic = "me.my"

type MyEvent struct {
    banji.EventEmbed
    someStuff string
}

type (e *MyEvent) Topic() string {
    return MyTopic
}

type (e *MyEvent) SomeStuff() string {
    return e.someStuff
}

// It may be helpful to define an initializer,
// especially if you want other components to create your event.

func NewMyEvent(someStuff string) *MyEvent {
    return &MyEvent{
        someStuff: someStuff,
    }
}

```

Likewise, receivers implement `banji.Receiver` by embedding `banji.ReceiverEmbed`.

```go

// In the same package or another package...

type MyReceiver struct {
    banji.ReceiverEmbed
    someOtherStuff string
}

type (r *MyReceiver) Topic() string {
    return MyTopic
}

type (r *MyReceiver) Handle(e Event) error {
    // You should not need to worry about validating type assertion
    // as long as you have namespaced your topic string (to prevent
    // conflicts) and aren't doing something weird.

    event, _ := e.(*MyEvent)
    fmt.Println(event.SomeStuff())
    fmt.Println(r.someOtherStuff)

    return nil
}

func NewMyReceiver(someOtherStuff string) *MyReceiver {
    return &MyReceiver{
        someOtherStuff: someOtherStuff,
    }
}

```
