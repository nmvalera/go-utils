## App Flags

The application supports the following configuration flags:

### Main Server Configuration

| Flag | Environment Variable | Description | Default Value |
|------|---------------------|-------------|---------------|
| `--main-ep-addr` | `MAIN_ENTRYPOINT_ADDR` | Main entrypoint address to listen on | `:8080` |
| `--main-ep-net-keep-alive` | `MAIN_ENTRYPOINT_NET_KEEP_ALIVE` | Main entrypoint keep alive | `0` |
| `--main-ep-http-read-timeout` | `MAIN_ENTRYPOINT_HTTP_READ_TIMEOUT` | Main entrypoint maximum duration for reading an entire request including the body (zero means no timeout) | `30s` |
| `--main-ep-http-read-header-timeout` | `MAIN_READ_HEADER_TIMEOUT` | Main entrypoint maximum duration for reading request headers (zero uses the value of read timeout) | `30s` |
| `--main-ep-http-write-timeout` | `MAIN_ENTRYPOINT_HTTP_WRITE_TIMEOUT` | Main entrypoint maximum duration for writing the response (zero means no timeout) | `30s` |
| `--main-ep-http-idle-timeout` | `MAIN_ENTRYPOINT_HTTP_IDLE_TIMEOUT` | Main entrypoint maximum amount of time to wait for the next request when keep-alives are enabled (zero uses the value of read timeout) | `30s` |

### Health Check Server Configuration

| Flag | Environment Variable | Description | Default Value |
|------|---------------------|-------------|---------------|
| `--healthz-ep-addr` | `HEALTHZ_ENTRYPOINT_ADDR` | Health check entrypoint address | `:8081` |
| `--healthz-ep-keep-alive` | `HEALTHZ_ENTRYPOINT_NET_KEEP_ALIVE` | Health check entrypoint keep alive | `0` |
| `--healthz-ep-http-read-timeout` | `HEALTHZ_ENTRYPOINT_HTTP_READ_TIMEOUT` | Health entrypoint maximum duration for reading an entire request including the body (zero means no timeout) | `30s` |
| `--healthz-ep-http-read-header-timeout` | `HEALTHZ_ENTRYPOINT_HTTP_READ_HEADER_TIMEOUT` | Health entrypoint maximum duration for reading request headers (zero uses the value of read timeout) | `30s` |
| `--healthz-ep-http-write-timeout` | `HEALTHZ_ENTRYPOINT_HTTP_WRITE_TIMEOUT` | Health entrypoint maximum duration for writing the response (zero means no timeout) | `30s` |
| `--healthz-ep-http-idle-timeout` | `HEALTHZ_ENTRYPOINT_HTTP_IDLE_TIMEOUT` | Health entrypoint maximum amount of time to wait for the next request when keep-alives are enabled (zero uses the value of read timeout) | `30s` |

All configuration flags can be set via command-line arguments, environment variables, or through a configuration file. The application uses [Viper](https://github.com/spf13/viper) for configuration management, which supports multiple configuration formats.
