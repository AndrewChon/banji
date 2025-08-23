# Banji

[![Scorecard supply-chain security](https://github.com/AndrewChon/banji/actions/workflows/scorecard.yml/badge.svg)](https://github.com/AndrewChon/banji/actions/workflows/scorecard.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/AndrewChon/banji)](https://goreportcard.com/report/github.com/AndrewChon/banji)
[![GoDoc](https://godoc.org/github.com/AndrewChon/banji?status.png)](https://godoc.org/github.com/AndrewChon/banji)

Banji is a concurrent event engine and EDA framework for Go.

> [!CAUTION]
> Banji is in the very, _very_ early stages of development.
> \
> There is undoubtedly room for improvement, and this codebase is bound to change. Please keep this in mind if you
> intend to use this package.

## About

I created Banji mainly for myself. The goal of Banji is to provide a structured and extensible event engine that can be
used for anything. There are some events we do not need to handle or care about; some capabilities we simply do not
need. Conversely, there are some components we may desire in the future but cannot reasonably foresee needing now.
Therefore, we shouldn't need to implement an interface and change our implementation every time some new component or
capability comes along.

## Features

- Concurrent demultiplexing
- Event priorities
- Built-in events for happenings within the engine
- Proactor-esque design (see `ErrorEvent`)
- Reflection-free routing
- Graceful shutdown mechanics
- Tick-based engine loop
- Unified logging system

## Quick Start

Check the [wiki](https://github.com/AndrewChon/banji/wiki) for a quick tutorial on how to create events, receivers, and
components.

```go
func main() {
    // This allows us to check for the program exiting.
    stop := make(chan os.Signal)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    eng := banji.New(
        banji.WithTPS(128),
        banji.WithDemuxers(8),
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