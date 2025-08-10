package banji

import (
	"runtime"
)

type Option func(*Options)

type Options struct {
	TPS      int
	Demuxers int
}

func NewOptions(opts ...Option) *Options {
	// Default settings.
	options := &Options{
		TPS:      128,
		Demuxers: runtime.NumCPU(),
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func WithTPS(n int) Option {
	if n < 1 {
		n = 1
	}

	return func(options *Options) {
		options.TPS = n
	}
}

func WithDemuxers(n int) Option {
	if n < 1 {
		n = 1
	}

	return func(options *Options) {
		options.Demuxers = n
	}
}
