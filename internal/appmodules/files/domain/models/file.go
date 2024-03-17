package models

import "bytes"

type File struct {
	ID     int
	UserId int
	Type   string
	Path   int
	Buffer bytes.Buffer
}
