package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func createZip(zipFileName string, folderToZip string) error {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	err = filepath.Walk(folderToZip, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(folderToZip, file)
		if err != nil {
			return err
		}

		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		fileToZip, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileToZip.Close()
		_, err = io.Copy(zipEntry, fileToZip)
		return err
	})

	return err
}

func unzip(zipFileName string, destFolder string) error {
	zipFile, err := zip.OpenReader(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	err = os.Mkdir(destFolder, 0755)
	_check(err)

	for _, file := range zipFile.File {
		filePath := filepath.Join(destFolder, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		destFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer destFile.Close()

		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		_, err = io.Copy(destFile, rc)
		if err != nil {
			return err
		}
	}

	return nil
}
