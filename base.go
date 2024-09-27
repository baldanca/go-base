// TODO: Add tests
// TODO: Add examples
// TODO: Add benchmarks
// TODO: Add README.md
// TODO: Makefile
// TODO: godoc

// Package base provides a way to initialize default resources that is used by almost all services.
// This package is used to load:
// - Context
// - Environment Variables
// - Logger
// - Time Resources
// - HTTP Client
package base

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog"
)

type (
	// Baser is the interface that wraps the basic methods to get the initialized resources.
	Baser[envType any] interface {
		Ctx() context.Context
		CancelCtx() func()
		Env() envType
		Logger() *zerolog.Logger
		TimeLocation() *time.Location
		// TimeNow returns the current time in the time location.
		TimeNow() time.Time
		HTTPClient() *http.Client
	}

	base[envType any] struct {
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
func New[envType any](cfg Config) Baser[envType] {
	cfg.validate()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	baser := base[envType]{
		ctx:       ctx,
		cancelCtx: cancel,
		logger:    cfg.Logger,
	}

	baser.loadEnviromentVariables()
	baser.loadTimeLocation(cfg)
	baser.loadHTTPClient(cfg)

	return baser
}

func (baser base[any]) loadTimeLocation(cfg Config) base[any] {
	location, err := time.LoadLocation(cfg.TimeLocation)
	if err != nil {
		panic("failed to load time location: " + err.Error())
	}

	baser.timeLocation = location

	return baser
}

func (baser base[any]) loadEnviromentVariables() base[any] {
	err := env.Parse(&baser.env)
	if err != nil {
		panic("failed to load env: " + err.Error())
	}

	return baser
}

func (baser base[any]) loadHTTPClient(cfg Config) base[any] {
	baser.httpClient = &http.Client{
		Transport: cfg.HTTPClientConfig.Transport,
		Jar:       cfg.HTTPClientConfig.Jar,
		Timeout:   cfg.HTTPClientConfig.Timeout,
	}

	return baser
}

func (baser base[any]) Ctx() context.Context {
	return baser.ctx
}

func (baser base[any]) CancelCtx() func() {
	return baser.cancelCtx
}

func (baser base[any]) Env() any {
	return baser.env
}

func (baser base[any]) Logger() *zerolog.Logger {
	return baser.logger
}

func (baser base[any]) TimeLocation() *time.Location {
	return baser.timeLocation
}

func (baser base[any]) TimeNow() time.Time {
	return time.Now().In(baser.timeLocation)
}

func (baser base[any]) HTTPClient() *http.Client {
	return baser.httpClient
}
