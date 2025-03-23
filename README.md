# Go Monorepo 

This monorepo leverages centralized documentation and atomic commits,
simplifying management across multiple applications and shared libraries.
Clear module boundaries and explicit dependency management ensure robust and maintainable code. 

## Structure
* `apps/`: Executable applications and services.
* `pkg/`: Shared libraries, organized by their functional responsibility.
* `mocks/`: Auto-generated mocks managed through `.mockery.yaml` for effective testing.
* `cmd/`: Code generators and other executable tools.

Each component is organized into its own Go module,
simplifying dependency management and ensuring clean boundaries between services.

## Architecture
### Go Modules
* Each app and package is encapsulated as a separate Go module.
* Root-level `go.work` file enables streamlined local development and inter-module dependencies.
* Internal packages (`internal/`) safeguard application-specific implementations from unintended external usage.

### Dependency Management
* Clear separation enforced by Go Modules.
* No cross-module import without explicit declaration, ensuring modularity and clarity.
* Internal packages (`internal/`) protect application-specific logic from unintended external dependencies.

## Modules
### Apps (`apps/`)
* `errdoc`: HTTP service to document and expose defined error codes.
* `examples`: Demonstrates practical usage of shared libraries.

### Shared Libraries (`pkg/`)
* `logger`: Structured logging abstraction supporting [zerolog](https://github.com/rs/zerolog).
* `alerter`: Notification system integration ([Slack](https://github.com/slack-go/slack)) with rate limiting.
* `errcode`: Centralized error codes definition and generation.
* `middleware`: Common HTTP middleware utilities.
* `fflag`: Feature flags using [OpenFeature SDK](https://github.com/open-feature/go-sdk).
* `environment`: Consistent environment configuration handling.

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