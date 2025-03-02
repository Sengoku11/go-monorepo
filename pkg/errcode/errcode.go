// Package errcode provides a centralized definition of error codes for the project,
// taking full advantage of the monorepo structure to ensure consistency across services.
//
// NOTE: The generator checks for duplicate code values across all constants of type Code.
// If two constants have the same numeric value, generation fails.
//
//go:generate go run ../../cmd/generrors/main.go
package errcode

// Code represents an application specific error code.
//
// NOTE: Avoid iota to keep numeric values stable. This helps ensure older references
// remain valid if they are sorted or new codes get added in the middle.
//
//go:generate stringer -type=Code -output=generated_errors_string.go
type Code int

// Generic error indicating failure in undefined upstream service.
const (
	ErrUpstream                   Code = 1000
	ErrUpstreamTimeout            Code = 1001
	ErrUpstreamInternal           Code = 1002
	ErrUpstreamNotFound           Code = 1003
	ErrUpstreamBadRequest         Code = 1004
	ErrUpstreamServiceUnavailable Code = 1005
	ErrUpstreamRateLimited        Code = 1006
	ErrUpstreamConnectionRefused  Code = 1007
	ErrUpstreamMalformedResponse  Code = 1008
	ErrUpstreamAuthFailed         Code = 1009
)
