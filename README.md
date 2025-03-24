# Go Monorepo 

This monorepo leverages centralized documentation and atomic commits,
simplifying management across multiple microservices.

## Structure
* `apps/`: Executable applications and services.
* `pkg/`: Shared libraries, organized by their functional responsibility.
* `mocks/`: Auto-generated mocks managed through `.mockery.yaml` for effective testing.
* `cmd/`: Code generators and other executable tools.

## Components
### Apps (`apps/`)
* `errdoc`: HTTP service to document and expose defined error codes.
* `examples/`: Demonstrates practical usage of shared libraries.
  * `canaryrollout`: Implements canary rollout with feature flags.
  * `killswitch`: Remotely turn on/off specific services within the app.

### Shared Libraries (`pkg/`)
* `logger`: Structured logging abstraction supporting [zerolog](https://github.com/rs/zerolog).
* `alerter`: Notification system integration ([Slack](https://github.com/slack-go/slack)) with rate limiting.
* `errcode`: Centralized error codes definition and generation.
* `middleware`: Common HTTP middleware utilities.
* `fflag`: Feature flags using [OpenFeature SDK](https://github.com/open-feature/go-sdk).
* `environment`: Consistent environment configuration handling.

## Module Organization
* Each app and package is encapsulated as a separate Go module.
* A root-level `go.work` file enables streamlined local development and inter-module dependencies.
* Internal packages (`internal/`) protect application-specific logic from unintended external dependencies.

## Development
### Prerequisites
* Go 1.24 or newer

### Installation Commands
```bash
# mockery
go install github.com/vektra/mockery/v2@v2.52.4

# stringer
go install golang.org/x/tools/cmd/stringer@v0.30.0

# documentation
go install golang.org/x/pkgsite/cmd/pkgsite@latest
```

### Running the Documentation Server
```bash
pkgsite -http :8080
```