package errors

import "errors"

var (
	ExternalAPIError = errors.New("external API error")
	MappingError     = errors.New("mapping error")
)
