package csv

import (
	"bytes"
	"encoding/csv"
	"sort"

	"golang.org/x/exp/maps"
)

type mapWriter struct {
	writer     *csv.Writer
	iobuffer   *bytes.Buffer
	allHeaders []string
}

func NewMapWriter() Writer {
	var csvBuffer bytes.Buffer
	return &mapWriter{
		iobuffer: &csvBuffer,
		writer:   csv.NewWriter(&csvBuffer),
	}
}

func (m *mapWriter) WriteHeaders(headers []string) error {
	sort.Strings(headers)
	m.allHeaders = headers
	return m.writer.Write(headers)
}

func (m *mapWriter) WriteLine(value map[string]string) error {
	if value == nil {
		return NewError(ErrContract, "values cannot be empty")
	}
	if m.allHeaders != nil && len(value) != len(m.allHeaders) {
		return NewError(ErrContract, "headers length is different of values length")
	}
	if m.allHeaders == nil {
		m.WriteHeaders(maps.Keys(value))
	}

	return m.write(value)
}

func (m *mapWriter) WriteLines(values []map[string]string) error {
	if values == nil {
		return NewError(ErrContract, "values cannot be empty")
	}
	if m.allHeaders != nil && len(values[0]) != len(m.allHeaders) {
		return NewError(ErrContract, "headers length is different of values length")
	}
	if m.allHeaders == nil {
		m.WriteHeaders(maps.Keys(values[0]))
	}

	for _, value := range values {
		if err := m.write(value); err != nil {
			return err
		}
	}
	return nil
}

func (m *mapWriter) GetData() (*bytes.Buffer, error) {
	if err := m.Close(); err != nil {
		return nil, err
	}
	return m.iobuffer, nil
}

func (m *mapWriter) Close() error {
	m.writer.Flush()
	return m.writer.Error()
}

func (m *mapWriter) write(value map[string]string) error {
	toWrite := make([]string, 0, len(m.allHeaders))
	for _, header := range m.allHeaders {
		toWrite = append(toWrite, value[header])
	}
	return m.writer.Write(toWrite)
}
