package database

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/ghthor/database/datatype"
	"github.com/ziutek/mymysql/mysql"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func saveFilesTx(tx mysql.Transaction, files []datatype.FormFile, dir string) (sha1Names []string, err error) {
	sha1Names, err = saveFiles(files, dir)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, RollbackError{rollbackErr, err}
		} else {
			return nil, err
		}
	}
	return
}

func saveFiles(files []datatype.FormFile, dir string) (sha1Names []string, err error) {
	for _, file := range files {
		filename, err := saveFile(file, dir)
		if err != nil {
			return nil, err
		}
		sha1Names = append(sha1Names, filename)
	}
	return
}

func saveFile(formFile datatype.FormFile, dir string) (sha1Name string, err error) {
	file, header := formFile.File, formFile.Header

	h := sha1.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return
	}

	sha1Name = hex.EncodeToString(h.Sum(nil))

	// Append the extension
	parts := strings.Split(header.Filename, ".")
	sha1Name += "." + parts[len(parts)-1]

	imgFile, err := os.Create(filepath.Join(dir, sha1Name))
	if err != nil {
		return
	}
	defer imgFile.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		return
	}
	_, err = io.Copy(imgFile, file)
	if err != nil {
		return
	}

	return sha1Name, nil
}
