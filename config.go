package ps2

type config struct {
	outputType outputType
}

func NewConfig(options []Option) *config {
	c := &config{}
	for _, opt := range options {
		opt(c)
	}
	return c
}

type outputType string

const (
	outputTypeJSON outputType = "json"
	outputTypeYAML outputType = "yaml"
	outputTypeTOML outputType = "toml"
)

type Option func(*config)

func WithOutputTypeJSON() Option {
	return func(c *config) {
		c.outputType = outputTypeJSON
	}
}

func WithOutputTypeYAML() Option {
	return func(c *config) {
		c.outputType = outputTypeYAML
	}
}

func WithOutputTypeTOML() Option {
	return func(c *config) {
		c.outputType = outputTypeTOML
	}
}
