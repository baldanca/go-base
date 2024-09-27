// TODO: Add tests
// TODO: Add examples
// TODO: Add benchmarks
// TODO: Add README.md
// TODO: Makefile
// TODO: godoc

// Package init provides a way to initialize default resources that is used by almost all services.
// This package is used to load the environment variables, logger, time location and time now.
package init

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog"
)

type (
	// Initer is the interface that wraps the basic methods to get the initialized resources.
	Initer[envType any] interface {
		Ctx() context.Context
		CancelCtx() func()
		Env() envType
		Logger() *zerolog.Logger
		TimeLocation() *time.Location
		// TimeNow returns the current time in the time location.
		TimeNow() time.Time
		HTTPClient() *http.Client
	}

	Init[envType any] struct {
		ctx          context.Context
		cancelCtx    func()
		env          envType
		logger       *zerolog.Logger
		timeLocation *time.Location
		httpClient   *http.Client
	}

	// Config is used to configure the resources that will be initialized.
	Config struct {
		// If not provided, a default logger will be created with the level set to Info and output to os.Stdout.
		Logger *zerolog.Logger
		// If not provided, the default location will be used (UTC).
		TimeLocation string
		// If not provided, the default configuration will be used.
		HTTPClientConfig HTTPClientConfig
	}

	HTTPClientConfig struct {
		// If not provided, the default transport will be used (http.DefaultTransport).
		Transport http.RoundTripper
		Jar       http.CookieJar
		// If not provided, the default timeout will be used (10 seconds).
		Timeout time.Duration
	}
)

func (cfg *Config) validate() {
	if cfg.Logger == nil {
		logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
		cfg.Logger = &logger
	}

	if cfg.TimeLocation == "" {
		cfg.TimeLocation = "UTC"
	}

	cfg.HTTPClientConfig.validate()
}

func (cfg *HTTPClientConfig) validate() {
	if cfg.Transport == nil {
		cfg.Transport = http.DefaultTransport
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}
}

// New initialize default resources that is used by almost all services.
// - Context
// - Environment Variables
// - Logger
// - Time Resources
// - HTTP Client
//
// To load the environment variables need to pass the envType as struct with the 'env' labels like the example below.
//
//	type envModel struct {
//		Environment string `env:"ENVIRONMENT"`
//		Version     string `env:"VERSION"`
//		Database    string `env:"DATABASE"`
//	}
func New[envType any](cfg Config) Initer[envType] {
	cfg.validate()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	initer := Init[envType]{
		ctx:       ctx,
		cancelCtx: cancel,
		logger:    cfg.Logger,
	}

	initer.loadEnviromentVariables()
	initer.loadTimeLocation(cfg)
	initer.loadHTTPClient(cfg)

	return initer
}

func (initer Init[any]) loadTimeLocation(cfg Config) Init[any] {
	location, err := time.LoadLocation(cfg.TimeLocation)
	if err != nil {
		panic("failed to load time location: " + err.Error())
	}

	initer.timeLocation = location

	return initer
}

func (initer Init[any]) loadEnviromentVariables() Init[any] {
	err := env.Parse(&initer.env)
	if err != nil {
		panic("failed to load env: " + err.Error())
	}

	return initer
}

func (initer Init[any]) loadHTTPClient(cfg Config) Init[any] {
	initer.httpClient = &http.Client{
		Transport: cfg.HTTPClientConfig.Transport,
		Jar:       cfg.HTTPClientConfig.Jar,
		Timeout:   cfg.HTTPClientConfig.Timeout,
	}

	return initer
}

func (initer Init[any]) Ctx() context.Context {
	return initer.ctx
}

func (initer Init[any]) CancelCtx() func() {
	return initer.cancelCtx
}

func (initer Init[any]) Env() any {
	return initer.env
}

func (initer Init[any]) Logger() *zerolog.Logger {
	return initer.logger
}

func (initer Init[any]) TimeLocation() *time.Location {
	return initer.timeLocation
}

func (initer Init[any]) TimeNow() time.Time {
	return time.Now().In(initer.timeLocation)
}

func (initer Init[any]) HTTPClient() *http.Client {
	return initer.httpClient
}
