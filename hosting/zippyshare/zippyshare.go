package zippyshare

import (
	"bytes"
	"context"
	"filehost/hosting"
	"filehost/lib"
	"io"
	"net/http"
	"regexp"
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
func (z *Service) Upload(ctx context.Context, filename string, filereader io.Reader) (url []hosting.URL, err error) {
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
