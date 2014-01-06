package database

import (
	"github.com/ghthor/database/datatype"
	"mime/multipart"
	"os"
	"path/filepath"
)

func testFile(filepathStr string) (datatype.FormFile, error) {
	file, err := os.Open(filepathStr)
	if err != nil {
		return datatype.FormFile{}, err
	}

	return datatype.FormFile{
		File: file,
		Header: &multipart.FileHeader{
			Filename: filepath.Base(filepathStr),
		},
	}, nil
}
