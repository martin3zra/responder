package respond

import "net/http"

//ErrorFormatter response format
type ErrorFormatter interface {
	// An integer HTTP Status code error
	Status() int
	// An integer coding the error type
	Code() int
	// A short localized string that describes the error
	Error() string
	//(optional) A long localized error description if needed.
	//It can contain precise information about which
	//parameter is missing, or what are the
	//acceptable values
	Description() *string
	// (optional) A URL to online documentation that provides
	//more information about the error
	InfoURL() *string
}

//ErrorDescriptor to be embedded
type ErrorDescriptor struct{}

//Status (optional) HTTP Status code.
//It can contain precise information about which
//HTTP Status code correspond the error
func (ErrorDescriptor) Status() int {
	return http.StatusInternalServerError
}

//Description (optional) A long localized error description if needed.
//It can contain precise information about which
//parameter is missing, or what are the
//acceptable values
func (ErrorDescriptor) Description() *string { return nil }

//InfoURL (optional) A URL to online documentation that provides
//more information about the error
func (ErrorDescriptor) InfoURL() *string { return nil }
