package errs

import "errors"

type APIError string

func (a APIError) Err() error {
	return errors.New(string(a))
}

func (a APIError) Str() string {
	return string(a)
}

const (
	UnauthorizedAPIError            APIError = "employee id is unauthorized"
	GenericMetricsAPIError          APIError = "could not retrieve metrics"
	FormatMismatchExpiresInAPIError APIError = "expires_in must be defined in minutes"
	SaveURLMappingFailedAPIError    APIError = "unable to save url mapping"
	DeleteURLFailedAPIError         APIError = "unable to save url mapping"
	URLMustBeProvidedAPIError       APIError = "a url must be given to be shortened"
)
