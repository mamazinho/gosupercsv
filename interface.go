package csv

import "bytes"

type Writer interface {
	WriteHeaders(headers []string) error
	WriteLine(values map[string]string) error
	WriteLines(values []map[string]string) error
	GetData() (*bytes.Buffer, error)
	Close() error
}
