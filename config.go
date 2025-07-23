package banji

type Configuration struct {
	TPS        int
	Demuxers   int
	Components []Component
}

type Option func(*Configuration)

func WithTPS(tps int) Option {
	if tps < 1 {
		tps = 1
	}

	return func(c *Configuration) {
		c.TPS = tps
	}
}

func WithDemuxers(demuxers int) Option {
	if demuxers < 1 {
		demuxers = 1
	}

	return func(c *Configuration) {
		c.Demuxers = demuxers
	}
}

func WithComponents(cs ...Component) Option {
	return func(c *Configuration) {
		c.Components = append(c.Components, cs...)
	}
}

func NewConfiguration(options ...Option) *Configuration {
	cfg := new(Configuration)
	for _, o := range options {
		o(cfg)
	}

	return cfg
}
