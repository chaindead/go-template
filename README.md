# Go Template: Must Read

It's not a framework or just another reinvented wheel, but essentially a comfortably configured combination of the zerolog and viper libraries

---

## Out of the Box We Get

> Write down what else might be needed

- [x] logger, config
- [x] server start with metrics
- [x] linter config
- [x] dockerfile
- [x] gitignore
- [x] graceful shutdown
- [x] git build info at application startup in logs
- [ ] Taskfile with basic commands (which ones are needed?)

## Usage

- All of the features listed below are located in `/internal` and are open for modification for your specific project
- The implementation of all these features is based on a small amount of code and is designed to be modified as needed
- When importing the template, you must replace all mentions of `github.com/chaindead/go-template`
  * `go.mod`
  * `golang-ci.yaml`

### Logger

We use the global zerolog without wrappers
```go
package app

import (
  "github.com/rs/zerolog/log"
)

func Any(){
  log.Info().Msg("Example msg") // the message should start with an uppercase letter
  //OUTPUT 14:55:19 INF internal/app/main.go:26 > Example msg

  log.Err(err).Send() // will include a message line and an error stack trace
  //OUTPUT 14:57:18 ERR cmd/app/main.go:48 > error="fake error" stack="main.go:37((*App).Shutdown); main.go:47(main)"
}
```

- The global logger is configured in `/internal/logger`
- It adds a log line (to all logs)
- It appends the error stack in a readable format (for proper error handling, see Error Handling)
- The linter forbids the use of the `log` package in favor of `github.com/rs/zerolog/log`

### Configuration

We use viper + pflag

```go
import (
	cfg "github.com/spf13/pflag"
)

var (
	setting = cfg.String("domain.setting", "default", "Example for cfg domain with name setting")
)
```

- The implementation is hidden in `/internal/config`
- Configuration is loaded from flags and YAML
- Priority of values:
  1. Default value from code
  2. Value from YAML
  3. Value from flag
- Checks for unknownValues (warn log)
- Direct use of viper in code is forbidden by the linter outside of the config package
- The linter forces the use of the alias `cfg` for the package

### Metrics

Again, we use global variables without wrappers, this time using the Prometheus and promauto packages

- Custom metrics are stored in one place: `internal/metrics`
- Metrics are registered on `prometheus.DefaultRegistry`, if an external provider is used

### Error Handling

Error handling is closely linked with logging, so we use the approach proposed by zerolog.  
Namely, we enforce the use of `pkg/errors` as a first-class citizen, all for preserving the error stack.

1. When creating an error, use `pkg/errors.New` instead of the standard `errors.New`.
2. When receiving an external error, be sure to wrap it using `errors.Wrap`.
3. When passing the error up, you can discard the wrappers—the stack will be printed when logging.
4. Additional information can be added via `fmt.Errorf` (do not use `errors.Wrapf`, to avoid losing the call stack).

- The linter ensures that these rules are followed:
  1. It forbids the `errors` package (in favor of `pkg/errors`)
  2. It checks that external errors are wrapped
  3. ⚠️ It does not check for double wrapping of an error, which would reset the stack at the moment of the second wrapping
- As a result, you get a stack trace as shown in the logger example

## What's Not Included Here?

- Reinventing wheels: we use popular open-source libraries for logging and configuration
- Pre-written modules: this is not a framework, but a foundation that makes it easier to start a project
- Complications/new approaches: the goal is to keep it simple, accessible, and explicit
