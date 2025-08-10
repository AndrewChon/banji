package bus

import (
	"runtime"
)

type Option func(*Options)

type Options struct {
	Demuxers     int
	ErrorBuilder func(error) Emittable
}

func NewOptions(opts ...Option) *Options {
	// Default settings.
	options := &Options{
		Demuxers: runtime.NumCPU(),
		ErrorBuilder: func(err error) Emittable {
			return nil
		},
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func WithDemuxers(n int) Option {
	if n < 1 {
		n = 1
	}

	return func(options *Options) {
		options.Demuxers = n
	}
}

func WithErrorBuilder(builder func(error) Emittable) Option {
	return func(options *Options) {
		options.ErrorBuilder = builder
	}
}
