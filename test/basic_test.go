package test

import (
	"testing"
	"time"

	"github.com/AndrewChon/banji"
)

const (
	TPS      = 128
	Demuxers = 8
)

func TestActivation(t *testing.T) {
	eng := banji.New(
		banji.WithTPS(TPS),
		banji.WithDemuxers(Demuxers),
	)

	eng.Start()
	defer eng.Stop()

	time.Sleep(1 * time.Second)
}

func TestBuiltins(t *testing.T) {
	eng := banji.New(
		banji.WithTPS(TPS),
		banji.WithDemuxers(Demuxers),
	)

	eng.Subscribe(new(StartReceiverTest))
	eng.Subscribe(new(StopReceiverTest))
	eng.Subscribe(new(PreTickReceiverTest))
	eng.Subscribe(new(PostTickReceiverTest))

	eng.Start()
	defer eng.Stop()

	time.Sleep(1 * time.Second)
}
