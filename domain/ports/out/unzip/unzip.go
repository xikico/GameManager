package unzip

import "errors"

var ErrPassword = errors.New("archive password required or invalid")

type UnzipFunc func(path string, password string) (string, error)
