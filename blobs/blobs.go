package blobs

import (
	"errors"
)

var ErrUnknown = errors.New("blob unknown")

type Blob struct {
	Name string
	Meta map[string]string
}
