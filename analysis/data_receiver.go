package analysis

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

//Проверка на то, что str является URL
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

//Получение тела ответа по URL
func getURLReadCloser(url string) (*io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return &resp.Body, nil
}

//Получение файла по данному пути
func getFileReadCloser(path string) (*io.ReadCloser, error) {
	resp, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	var reader io.ReadCloser = resp

	return &reader, nil
}

//Функция для вычитывания источников данных
func GetData(paths []string, dataChan chan<- PathReader) []error {
	dataErrors := make([]error, 0)
	var err error

	for _, elem := range paths {
		var reader *io.ReadCloser

		if IsUrl(elem) {
			reader, err = getURLReadCloser(elem)
		} else {
			reader, err = getFileReadCloser(elem)
		}
		if err != nil {
			dataErrors = append(dataErrors, fmt.Errorf("Failed to get data from resource: %v\n", elem))
			continue
		}
		dataChan <- PathReader{Path: elem, Reader: reader}
	}
	return dataErrors
}
