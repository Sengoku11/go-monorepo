# Go Monorepo 

A `go.work` template with prebuilt instruments.


### Features
* Automatic .env loader for local development.
* Recursive `make tidy` and `make lint`.
* Wrapped [zerolog](https://github.com/rs/zerolog) logger.
* Alerter hooks with Slack implementation and per-message rate limiting.
* Custom middlewares.
* Mock config to generate mocks for all packages.
* Error codes with auto-generation tools.

### Requirements
```bash
# mockery:
go install github.com/vektra/mockery/v2@v2.52.4
```
```bash
#stringer
go install golang.org/x/tools/cmd/stringer@v0.30.0
```