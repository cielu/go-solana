package core

import "errors"

// NotFound is returned by API methods if the requested item does not exist.
var NotFound = errors.New("not found")

