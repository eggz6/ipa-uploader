package uploader

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/eggz6/ipa-uploader/ipa"
)

type OSSUploader struct {
	cli      *oss.Client
	b        *oss.Bucket
	endpoint string
	bucket   string
}

func NewOss(bucket, endpoint, accessID, accessKey string) (ipa.Uploader, error) {
	client, err := oss.New(endpoint, accessID, accessKey)
	if err != nil {
		return nil, err
	}

	b, err := client.Bucket(bucket)
	if err != nil {
		return nil, err
	}

	return &OSSUploader{
		b:        b,
		cli:      client,
		bucket:   bucket,
		endpoint: endpoint,
	}, nil
}

func (o *OSSUploader) Upload(ctx context.Context, r io.Reader, path string) (url string, err error) {
	path = strings.TrimPrefix(path, "/")
	err = o.b.PutObject(path, r)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.%s/%s", o.bucket, o.endpoint, path), nil
}
