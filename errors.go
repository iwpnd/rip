package rip

import "errors"

// ErrInvalidRequestBody occurs when Request body is of unknown type
// and/or not mashalable.
var ErrInvalidRequestBody = errors.New("invalid request body")

// ErrClientMissing occurs when Request is instantiated without Client.NR()
var ErrClientMissing = errors.New("use .NR() to create a new request instead")
