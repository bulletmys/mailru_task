package analysis

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

//Проверка на то, что str является URL.
func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

//Получение тела ответа по URL.
func getURLReadCloser(url string) (*io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't get data from URL: %w", err)
	}

	return &resp.Body, nil
}

//Получение файла по данному пути.
func getFileReadCloser(path string) (*io.ReadCloser, error) {
	resp, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("can't open file: %w", err)
	}

	var reader io.ReadCloser = resp

	return &reader, nil
}

//Функция для вычитывания источников данных.
func GetData(paths []string, dataChan chan<- PathReader) []error {
	dataErrors := make([]error, 0)

	var err error

	for _, path := range paths {
		var reader *io.ReadCloser

		if IsURL(path) {
			reader, err = getURLReadCloser(path)
		} else {
			reader, err = getFileReadCloser(path)
		}

		if err != nil {
			dataErrors = append(dataErrors, fmt.Errorf("failed to get data from resource %v: %w", path, err))
		} else {
			dataChan <- PathReader{Path: path, Reader: reader}
		}
	}

	return dataErrors
}
