package hosting

import (
	"context"
	"io"
)

// URL is struc for url
type URL struct {
	StartByte int64
	EndByte   int64
	URI       string
	Expire    int
}

// Service is interface
type Service interface {
	Upload(ctx context.Context, filename string, filereader io.Reader) []URL
}
