package app

import (
	"github.com/kkrt-labs/go-utils/spf13"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	mainEntrypointFlag = &spf13.StringFlag{
		ViperKey:     "app.main-entrypoint.addr",
		Name:         "main-ep-addr",
		Env:          "MAIN_ENTRYPOINT_ADDR",
		Description:  "Main entrypoint address to listen on",
		DefaultValue: defaultConfig.MainEntrypoint.Addr,
	}
	mainKeepAliveFlag = &spf13.StringFlag{
		ViperKey:     "app.main-entrypoint.net.keep-alive",
		Name:         "main-ep-net-keep-alive",
		Env:          "MAIN_ENTRYPOINT_NET_KEEP_ALIVE",
		Description:  "Main entrypoint keep alive",
		DefaultValue: defaultConfig.MainEntrypoint.Net.KeepAlive,
	}
	mainReadTimeoutFlag = &spf13.StringFlag{
		ViperKey:     "app.main-entrypoint.http.read-timeout",
		Name:         "main-ep-http-read-timeout",
		Env:          "MAIN_ENTRYPOINT_HTTP_READ_TIMEOUT",
		Description:  "Main entrypoint maximum duration for reading an entire request including the body (zero means no timeout)",
		DefaultValue: defaultConfig.MainEntrypoint.HTTP.ReadTimeout,
	}
	mainReadHeaderTimeoutFlag = &spf13.StringFlag{
		ViperKey:     "app.main-entrypoint.http.read-header-timeout",
		Name:         "main-ep-http-read-header-timeout",
		Env:          "MAIN_ENTRYPOINT_HTTP_READ_HEADER_TIMEOUT",
		Description:  "Main entrypoint maximum duration for reading request headers (zero uses the value of read timeout)",
		DefaultValue: defaultConfig.MainEntrypoint.HTTP.ReadHeaderTimeout,
	}
	mainWriteTimeoutFlag = &spf13.StringFlag{
		ViperKey:     "app.main-entrypoint.http.write-timeout",
		Name:         "main-ep-http-write-timeout",
		Env:          "MAIN_ENTRYPOINT_HTTP_WRITE_TIMEOUT",
		Description:  "Main entrypoint maximum duration for writing the response (zero means no timeout)",
		DefaultValue: defaultConfig.MainEntrypoint.HTTP.WriteTimeout,
	}
	mainIdleTimeoutFlag = &spf13.StringFlag{
		ViperKey:     "app.main-entrypoint.http.idle-timeout",
		Name:         "main-ep-http-idle-timeout",
		Env:          "MAIN_ENTRYPOINT_HTTP_IDLE_TIMEOUT",
		Description:  "Main entrypoint maximum amount of time to wait for the next request when keep-alives are enabled (zero uses the value of read timeout)",
		DefaultValue: defaultConfig.MainEntrypoint.HTTP.IdleTimeout,
	}
	healthzEntrypointFlag = &spf13.StringFlag{
		ViperKey:     "app.healthz-entrypoint.addr",
		Name:         "healthz-ep-addr",
		Env:          "HEALTHZ_ENTRYPOINT_ADDR",
		Description:  "Healthz entrypoint address to listen on",
		DefaultValue: defaultConfig.HealthzEntrypoint.Addr,
	}
	healthzKeepAliveFlag = &spf13.StringFlag{
		ViperKey:     "app.healthz-entrypoint.net.keep-alive",
		Name:         "healthz-ep-net-keep-alive",
		Env:          "HEALTHZ_ENTRYPOINT_NET_KEEP_ALIVE",
		Description:  "Healthz entrypoint keep alive",
		DefaultValue: defaultConfig.HealthzEntrypoint.Net.KeepAlive,
	}
	healthzReadTimeoutFlag = &spf13.StringFlag{
		ViperKey:     "app.healthz-entrypoint.http.read-timeout",
		Name:         "healthz-ep-http-read-timeout",
		Env:          "HEALTHZ_ENTRYPOINT_HTTP_READ_TIMEOUT",
		Description:  "Healthz entrypoint maximum duration for reading an entire request including the body (zero means no timeout)",
		DefaultValue: defaultConfig.HealthzEntrypoint.HTTP.ReadTimeout,
	}
	healthzReadHeaderTimeoutFlag = &spf13.StringFlag{
		ViperKey:     "app.healthz-entrypoint.http.read-header-timeout",
		Name:         "healthz-ep-http-read-header-timeout",
		Env:          "HEALTHZ_ENTRYPOINT_HTTP_READ_HEADER_TIMEOUT",
		Description:  "Healthz entrypoint maximum duration for reading request headers (zero uses the value of read timeout)",
		DefaultValue: defaultConfig.HealthzEntrypoint.HTTP.ReadHeaderTimeout,
	}
	healthzWriteTimeoutFlag = &spf13.StringFlag{
		ViperKey:     "app.healthz-entrypoint.http.write-timeout",
		Name:         "healthz-ep-http-write-timeout",
		Env:          "HEALTHZ_ENTRYPOINT_HTTP_WRITE_TIMEOUT",
		Description:  "Healthz entrypoint maximum duration for writing the response (zero means no timeout)",
		DefaultValue: defaultConfig.HealthzEntrypoint.HTTP.WriteTimeout,
	}
	healthzIdleTimeoutFlag = &spf13.StringFlag{
		ViperKey:     "app.healthz-entrypoint.http.idle-timeout",
		Name:         "healthz-ep-http-idle-timeout",
		Env:          "HEALTHZ_ENTRYPOINT_HTTP_IDLE_TIMEOUT",
		Description:  "Healthz entrypoint maximum amount of time to wait for the next request when keep-alives are enabled (zero uses the value of read timeout)",
		DefaultValue: defaultConfig.HealthzEntrypoint.HTTP.IdleTimeout,
	}
)

func AddFlags(v *viper.Viper, f *pflag.FlagSet) {
	mainEntrypointFlag.Add(v, f)
	mainKeepAliveFlag.Add(v, f)
	mainReadTimeoutFlag.Add(v, f)
	mainReadHeaderTimeoutFlag.Add(v, f)
	mainWriteTimeoutFlag.Add(v, f)
	mainIdleTimeoutFlag.Add(v, f)
	healthzEntrypointFlag.Add(v, f)
	healthzKeepAliveFlag.Add(v, f)
	healthzReadTimeoutFlag.Add(v, f)
	healthzReadHeaderTimeoutFlag.Add(v, f)
	healthzWriteTimeoutFlag.Add(v, f)
	healthzIdleTimeoutFlag.Add(v, f)
}
