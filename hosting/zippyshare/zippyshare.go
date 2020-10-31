package zippyshare

import (
	"bytes"
	"context"
	"errors"
	"filehost/hosting"
	"filehost/lib"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// Service for zippyshare.com
type Service struct{}

const maxUpload = 500 * 1000 * 1000
const expire = 30 * 24 * time.Hour

type uploadSession struct {
	id     string
	server string
}

// Upload is func to upload a file
func (s *Service) Upload(ctx context.Context, filename string, filereader io.Reader) (url []hosting.URL, err error) {
	// prep http client
	client := http.DefaultClient

	// create upload session
	ses, err := createSession(client)
	if err != nil {
		return
	}

	// prep form
	form := lib.NewForm()
	form.AddField("uploadid", ses.id)
	form.AddField("notprivate", "false")
	form.AddField("zipname", "")
	form.AddField("ziphash", "")
	form.AddField("embPlayerValues", "false")
	position := int64(0)
	for {
		readed, err := form.AddFileField("file", filename, io.LimitReader(filereader, maxUpload))
		if err != nil {
			return []hosting.URL{}, err
		}
		url = append(url, hosting.URL{
			StartByte: position,
			EndByte:   position + readed,
		})
		position = position + readed
		if readed < maxUpload {
			break
		}
	}
	form.Close()

	// uploading
	res, err := client.Post(ses.server, form.ContentType(), form.Buff)
	if err != nil {
		return
	}
	uri := fetchDownloadURI(res)
	for k, v := range uri {
		url[k].URI = v
		url[k].Expire = time.Now().Add(expire).Nanosecond()
	}
	return
}

// DirectDownloadURI is func to generate direct download link
func (s *Service) DirectDownloadURI(ctx context.Context, uri string) (duri string, err error) {
	re := regexp.MustCompile("\\/\\/www([\\d]+)\\.zippyshare\\.com\\/v\\/([^\\/]+)")
	match := re.FindStringSubmatch(uri)
	if len(match) != 3 {
		err = errors.New("Invalid uri")
		return "", err
	}
	duri = "https://www" + match[1] + ".zippyshare.com/d/" + match[2] + "/"

	client := http.DefaultClient
	res, err := client.Get(uri)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	body := buf.String()

	re = regexp.MustCompile("var a = ([\\d]+)\\;")
	match = re.FindStringSubmatch(body)
	if len(match) != 2 {
		err = errors.New("Unable to get unique download code")
		return "", err
	}
	a, _ := strconv.Atoi(match[1])
	duri = duri + strconv.Itoa((a*a*a)+3) + "/"

	re = regexp.MustCompile("\\(Math\\.pow\\(a, 3\\)\\+b\\)\\+\\\"\\/([^\\\"]+)\\\";")
	match = re.FindStringSubmatch(body)
	if len(match) != 2 {
		err = errors.New("Unable to get filename")
		return "", err
	}
	duri = duri + match[1]
	return duri, nil
}

func createSession(client *http.Client) (ses *uploadSession, err error) {
	res, err := client.Get("https://www.zippyshare.com")
	if err != nil {
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	body := buf.String()

	re := regexp.MustCompile("var uploadId = \\'([\\w\\d]+)\\'\\;")
	uploadID := re.FindStringSubmatch(body)[1]
	re = regexp.MustCompile("var server = \\'([\\w\\d]+)\\'\\;")
	server := "https://" + re.FindStringSubmatch(body)[1] + ".zippyshare.com/upload"

	ses = &uploadSession{
		id:     uploadID,
		server: server,
	}
	return
}

func fetchDownloadURI(res *http.Response) (urls []string) {
	var buff bytes.Buffer
	buff.ReadFrom(res.Body)
	body := buff.String()

	// print(body)
	re := regexp.MustCompile("\\[url\\=([^\\]]+)\\]")
	match := re.FindAllStringSubmatch(body, -1)
	for _, v := range match {
		urls = append(urls, v[1])
	}
	return
}
