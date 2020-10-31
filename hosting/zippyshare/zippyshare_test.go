package zippyshare_test

import (
	"context"
	"filehost/hosting/zippyshare"
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
