package utils

import "errors"

var (
	PublicIPEmpty        error = errors.New("public ip can not be empty")
	UnSupportMultiRecord error = errors.New("unsupported multiple records")
	NoChangeSkip         error = errors.New("no change, Skip")
)
