package datatype

import (
	"io"
	"mime/multipart"
)

type (
	Id uint64
)

type UploadedTempFile interface {
	io.Reader
	io.Seeker
	io.Closer
}

type FormFile struct {
	File   UploadedTempFile
	Header *multipart.FileHeader
}
