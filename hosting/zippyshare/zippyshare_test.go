package zippyshare_test

import (
	"context"
	"github.com/azamaulanaaa/go-filehost/hosting/zippyshare"
	"os"
	"testing"
)

const filename = "zippyshare_test.go"

func TestUpload(t *testing.T) {
	file, err := os.Open(filename)
	if err != nil {
		t.Error(err)
	}
	zippy := &zippyshare.Service{}
	_, err = zippy.Upload(context.TODO(), filename, file)
	if err != nil {
		t.Error(err)
	}
}

func TestDirectDownloadURI(t *testing.T) {
	zippy := &zippyshare.Service{}
	_, err := zippy.DirectDownloadURI(context.TODO(), "https://www106.zippyshare.com/v/1SMt0BDP/file.html")
	if err != nil {
		t.Error(err)
	}
}
