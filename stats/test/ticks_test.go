package hello_world

import (
	"banji"
	"banji/stats"
	"fmt"
	"testing"
	"time"
)

func TestTicks(_ *testing.T) {
	s := stats.New()

	eng := banji.New(
		banji.WithTPS(128),
		banji.WithDemuxers(8),
		banji.WithComponents(s),
	)

	eng.Start()
	defer eng.Stop()

	time.Sleep(time.Second * 3)

	fmt.Printf("%.2f\n", s.Ticks().PerSecond())
}
