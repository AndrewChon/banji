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

### "Hello, world!" Example

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndrewChon/banji"
)

// We define a receiver to listen for banji.StartEvent, which has the topic string banji.StartTopic.
type Greeter struct {
	banji.ReceiverEmbed
}

func (g *Greeter) Topic() string {
	return banji.StartTopic
}

func (g *Greeter) Handle(event banji.Event) error { 
	// For our purposes, we don't need any data or fields from the event.
	_, _ = event.(*banji.StartEvent)
	
	fmt.Println("Hello, world!")

	return nil
}

// Terminal output: "Hello, world!"
func main() {
	// This allows us to check for the program exiting.
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	eng := banji.New(
		banji.WithTPS(128),
		banji.WithDemuxers(8),
	)

	// We register our Greeter receiver to the engine.
	eng.Subscribe(new(Greeter))

	eng.Start() // This is non-blocking.
	defer eng.Stop()

	<-stop
}

```