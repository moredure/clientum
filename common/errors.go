package common

import "errors"

var (
	ServerAddressErr = errors.New("$SERVER_URL env is missing")
	UserMissingErr   = errors.New("$USER env is missing")
)
