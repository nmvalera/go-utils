package http

import (
	"net"
	"net/http"
	"time"

	"github.com/nmvalera/go-utils/common"
	"github.com/nmvalera/go-utils/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// DefaultEntrypointConfig returns a default EntrypointConfig.
func DefaultEntrypointConfig() *EntrypointConfig {
	return &EntrypointConfig{
		HTTP: &ServerConfig{
			ReadTimeout:       common.Ptr(30 * time.Second),
			ReadHeaderTimeout: common.Ptr(30 * time.Second),
			WriteTimeout:      common.Ptr(30 * time.Second),
			IdleTimeout:       common.Ptr(30 * time.Second),
			MaxHeaderBytes:    common.Ptr(http.DefaultMaxHeaderBytes),
		},
		Net: &ListenConfig{
			KeepAlive: common.Ptr(-time.Second),
			KeepAliveProbe: &KeepAliveProbeConfig{
				Enable:   common.Ptr(false),
				Idle:     common.Ptr(15 * time.Second),
				Interval: common.Ptr(15 * time.Second),
				Count:    common.Ptr(9),
			},
		},
		TLS: &TLSCertConfig{},
	}
}

// EntrypointConfig is the configuration for an entrypoint.
type EntrypointConfig struct {
	Addr *string        `key:"addr,omitempty" desc:"TCP Address to listen on"`
	HTTP *ServerConfig  `key:"http,omitempty"`
	Net  *ListenConfig  `key:"net,omitempty"`
	TLS  *TLSCertConfig `key:"tls,omitempty"`
}

func (cfg *EntrypointConfig) Entrypoint() (*Entrypoint, error) {
	return NewEntrypoint(common.Val(cfg.Addr), WithServer(cfg.HTTP.Server()), WithListenConfig(cfg.Net.ListenConfig()))
}

type TLSCertConfig struct {
	CertFile *string `key:"cert-file,omitempty" env:"CERT_FILE" desc:"Path to the certificate file"`
	KeyFile  *string `key:"key-file,omitempty" env:"KEY_FILE" desc:"Path to the key file"`
}

type embedConfig struct {
	EP *EntrypointConfig `key:"ep,omitempty"`
}

// Env returns the environment variables for the entrypoint config.
// All environment variables are prefixed with "EP_".
func (cfg *EntrypointConfig) Env() (map[string]string, error) {
	return config.Env(&embedConfig{cfg}, nil)
}

// Unmarshal unmarshals the given viper into the entrypoint config.
// Assumes
// - all viper keys are prefixed with "ep."
// - all environment variables are prefixed with "EP_".
func (cfg *EntrypointConfig) Unmarshal(v *viper.Viper) error {
	return config.Unmarshal(&embedConfig{cfg}, v)
}

// AddFlags adds flags to the given viper and pflag.FlagSet.
// Sets
// - all viper keys with "ep." prefix
// - all environment variables with "EP_" prefix
// - all flags with "ep-" prefix
func AddFlags(v *viper.Viper, f *pflag.FlagSet) error {
	return config.AddFlags(&embedConfig{DefaultEntrypointConfig()}, v, f, nil)
}

type ServerConfig struct {
	ReadTimeout       *time.Duration `key:"read-timeout,omitempty" env:"READ_TIMEOUT" flag:"read-timeout" desc:"Maximum duration for reading the entire request including the body (zero means no timeout)"`
	ReadHeaderTimeout *time.Duration `key:"read-header-timeout,omitempty" env:"READ_HEADER_TIMEOUT" flag:"read-header-timeout" desc:"Maximum duration for reading request headers (zero uses the value of read timeout)"`
	WriteTimeout      *time.Duration `key:"write-timeout,omitempty" env:"WRITE_TIMEOUT" flag:"write-timeout" desc:"Maximum duration before timing out writes of the response (zero means no timeout)"`
	IdleTimeout       *time.Duration `key:"idle-timeout,omitempty" env:"IDLE_TIMEOUT" flag:"idle-timeout" desc:"Maximum duration to wait for the next request when keep-alives are enabled (zero uses the value of read timeout)"`
	MaxHeaderBytes    *int           `key:"max-header-bytes,omitempty" env:"MAX_HEADER_BYTES" flag:"max-header-bytes" desc:"Maximum number of bytes the server will read parsing the request header's keys and values, including the request line"`
}

func (cfg *ServerConfig) Server() *http.Server {
	return &http.Server{
		MaxHeaderBytes:    common.Val(cfg.MaxHeaderBytes),
		ReadTimeout:       common.Val(cfg.ReadTimeout),
		ReadHeaderTimeout: common.Val(cfg.ReadHeaderTimeout),
		WriteTimeout:      common.Val(cfg.WriteTimeout),
		IdleTimeout:       common.Val(cfg.IdleTimeout),
	}
}

type ListenConfig struct {
	KeepAlive      *time.Duration        `key:"keep-alive,omitempty" env:"KEEP_ALIVE" flag:"keep-alive" desc:"Keep alive period for network connections accepted by this entrypoint"`
	KeepAliveProbe *KeepAliveProbeConfig `key:"keep-alive-probe,omitempty" env:"KEEP_ALIVE_PROBE" flag:"keep-alive-probe"`
}

func (cfg *ListenConfig) ListenConfig() *net.ListenConfig {
	netCfg := &net.ListenConfig{
		KeepAlive: common.Val(cfg.KeepAlive),
	}

	if cfg.KeepAliveProbe != nil {
		netCfg.KeepAliveConfig = *cfg.KeepAliveProbe.KeepAliveProbe()
	}

	return netCfg
}

type KeepAliveProbeConfig struct {
	Enable   *bool          `key:"enable,omitempty" desc:"Enable keep alive probes"`
	Idle     *time.Duration `key:"idle,omitempty" desc:"Time that the connection must be idle before the first keep-alive probe is sent"`
	Interval *time.Duration `key:"interval,omitempty" desc:"Time between keep-alive probes"`
	Count    *int           `key:"count,omitempty" desc:"Maximum number of keep-alive probes that can go unanswered before dropping a connection"`
}

func (cfg *KeepAliveProbeConfig) KeepAliveProbe() *net.KeepAliveConfig {
	return &net.KeepAliveConfig{
		Enable:   common.Val(cfg.Enable),
		Count:    common.Val(cfg.Count),
		Idle:     common.Val(cfg.Idle),
		Interval: common.Val(cfg.Interval),
	}
}
