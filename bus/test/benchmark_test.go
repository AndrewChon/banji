package test

import (
	"testing"

	"github.com/AndrewChon/banji/bus"
)

const (
	Demuxers = 8
)

// BenchmarkPost serves to benchmark the underlying data structure of bus.Bus; more specifically, the efficiency of
// posting.
func BenchmarkPost(b *testing.B) {
	bs := bus.NewBus[*MockEmittable, *MockSubscriber[*MockEmittable]](
		bus.WithDemuxers(Demuxers),
	)

	for b.Loop() {
		bs.Post(new(MockEmittable), 0)
	}

	bs.Tick() // Drain
}

// BenchmarkTick serves to benchmark the underlying data structure of bus.Bus; more specifically, the efficiency of the
// drain operation and iterating over the working queue. The time per operation reported represents the total time it
// takes to run a tick divided by the number of events within the bus.
func BenchmarkTick(b *testing.B) {
	b.ResetTimer()
	bs := bus.NewBus[*MockEmittable, *MockSubscriber[*MockEmittable]](
		bus.WithDemuxers(Demuxers),
	)

	for b.Loop() {
		bs.Post(new(MockEmittable), 0)
	}

	b.ResetTimer()
	b.StartTimer()
	bs.Tick()
	b.StopTimer()

	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

// BenchmarkParallelPost serves to benchmark the underlying data structure of bus.Bus; more specifically, the efficiency
// of posting in a parallel setting. Due to the bus's concurrency-safety mechanisms, a slower time per operation is to
// be expected.
func BenchmarkParallelPost(b *testing.B) {
	bs := bus.NewBus[*MockEmittable, *MockSubscriber[*MockEmittable]](
		bus.WithDemuxers(Demuxers),
	)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bs.Post(new(MockEmittable), 0)
		}
	})

	bs.Tick() // Drain
}
