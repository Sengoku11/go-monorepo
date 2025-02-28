# Go Monorepo 

A `go.work` template with prebuilt instruments.


### Features
* Automatic .env loader for local development.
* Recursive `make tidy` and `make lint`.
* Wrapped [zerolog](https://github.com/rs/zerolog) logger.
* Alerter hooks with Slack implementation and per-message rate limiting.
* Custom middlewares.
* Mock config to generate mocks for all packages.

### Requirements
```bash
# mockery:
go install github.com/vektra/mockery/v2@v2.52.4
```