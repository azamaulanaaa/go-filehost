package lib

import (
	"bytes"
	"io"
	"mime/multipart"
	"strings"
)

// Form for post form
type Form struct {
	Buff        *bytes.Buffer
	fieldWriter *multipart.Writer
}

// NewForm is func to create a new fields
func NewForm() *Form {
	fields := &Form{
		Buff: &bytes.Buffer{},
	}
	fields.fieldWriter = multipart.NewWriter(fields.Buff)
	return fields
}

// AddField is func to add an ordinary field
func (f *Form) AddField(key string, value string) (err error) {
	var writer io.Writer

	writer, err = f.fieldWriter.CreateFormField(key)
	if err != nil {
		return
	}

	_, err = io.Copy(writer, strings.NewReader(value))
	if err != nil {
		return
	}

	return
}

// AddFileField is func to add file field
func (f *Form) AddFileField(key string, filename string, filereader io.Reader) (n int64, err error) {
	var writer io.Writer

	writer, err = f.fieldWriter.CreateFormFile(key, filename)
	if err != nil {
		return
	}

	return io.Copy(writer, filereader)
}

// Close is func to close Form
func (f *Form) Close() {
	f.fieldWriter.Close()
}

// ContentType is form content type
func (f *Form) ContentType() string {
	return f.fieldWriter.FormDataContentType()
}
