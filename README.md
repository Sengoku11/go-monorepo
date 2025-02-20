# Go Monorepo 

A `go.work` template with prebuilt instruments.


### Features
* Automatic .env loader for local development.
* Recursive `make tidy` and `make lint`.
* Wrapped [zerolog](https://github.com/rs/zerolog) logger.
* Alerter hooks with Slack implementation and per-message rate limiting.
