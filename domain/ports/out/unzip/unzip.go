package unzip

type UnzipFunc func(path string, password string) (string, error)
