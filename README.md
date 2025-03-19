# Go Monorepo 

This repository centralizes documentation and streamlines changes across client and server codebases,
showcasing the inherent benefits of a monorepo structure—such as the ability
to perform atomic commits—to ensure consistency and ease of management.

By leveraging `go.work`, imports are simplified, and integration with private repositories is seamless.

### Features
* Automatic `.env` loader for local development, plus recursive Make commands like `make tidy` and `make lint`.
* Wrapped [Zerolog](https://github.com/rs/zerolog) logger
* Alerter hooks with [Slack](https://github.com/slack-go/slack) and per-message rate limiting.
* Custom middlewares.
* Feature flags via [OpenFeature](https://github.com/open-feature/go-sdk) and [Flipt](https://www.flipt.io/). 
* Auto-generated mocks via `.mockery.yaml`.
* Auto-generated documentation for all error codes.

### Running the Documentation Server
```bash
pkgsite -http :8080
```

### Requirements
```bash
# mockery
go install github.com/vektra/mockery/v2@v2.52.4

# stringer
go install golang.org/x/tools/cmd/stringer@v0.30.0

# documentation
go install golang.org/x/pkgsite/cmd/pkgsite@latest
```