package models

import "bytes"

type File struct {
	ID     int
	UserID int
	Type   string
	Path   int
	Buffer bytes.Buffer
}
