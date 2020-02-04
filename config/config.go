package config

import "time"

// Config struc holds all the application configurations and uses stuct tags to
// define the environment variables responsible for each parameter and default values
type Config struct {
	UpstreamTimeout time.Duration `env:"BURROW_UPSTREAM_TIMEOUT" envDefault:"2000ms"`
	UpstreamServer  string        `env:"BURROW_UPSTREAM_SERVER" envDefault:"1.1.1.1"`
	UpstreamPort    string        `env:"BURROW_UPSTREAM_PORT" envDefault:"853"`
	EnableTCP       bool          `env:"BURROW_ENABLE_TCP" envDefault:"true"`
	EnableUDP       bool          `env:"BURROW_ENABLE_UDP" envDefault:"true"`
}
