package cos

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Client struct {
	C *cos.Client
}

func NewClient(urlStr, SecretID, SecretKey string) Client {
	u, err := url.Parse(urlStr)
	if err != nil {
		hlog.Fatal(err)
	}
	b := &cos.BaseURL{
		BucketURL: u,
	}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  SecretID,
			SecretKey: SecretKey,
		}})
	return Client{C: client}
}

func (c *Client) Upload(ctx context.Context, key string, file *multipart.FileHeader) error {
	v, _, err := c.C.Object.InitiateMultipartUpload(ctx, key, nil)
	if err != nil {
		return err
	}
	f, err := file.Open()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	partNum, partSize := cos.DividePart(file.Size, 16)
	var chunk cos.Chunk
	var chunks []cos.Chunk
	for i := int64(0); i < partNum; i++ {
		chunk.Number = int(i + 1)
		chunk.OffSet = i * partSize
		chunk.Size = partSize
		chunks = append(chunks, chunk)
	}

	if file.Size%partSize > 0 {
		chunk.Number = len(chunks) + 1
		chunk.OffSet = int64(len(chunks)) * partSize
		chunk.Size = file.Size % partSize
		chunks = append(chunks, chunk)
		partNum++
	}
	optCom := &cos.CompleteMultipartUploadOptions{}
	var wg sync.WaitGroup
	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk cos.Chunk) {
			defer wg.Done()
			_, err := f.Seek(chunk.OffSet, io.SeekStart)
			if err != nil {
				hlog.CtxErrorf(ctx, "%v", err)
			}
			resp, err := c.C.Object.UploadPart(
				ctx,
				key,
				v.UploadID,
				chunk.Number,
				cos.LimitReadCloser(f, chunk.Size),
				&cos.ObjectUploadPartOptions{ContentLength: chunk.Size})
			if err != nil {
				hlog.CtxErrorf(ctx, "%v", err)
			}
			optCom.Parts = append(optCom.Parts, cos.Object{ETag: resp.Header.Get("ETag"), PartNumber: chunk.Number})
			time.Sleep(time.Millisecond)
		}(chunk)
	}
	wg.Wait()
	_, _, err = c.C.Object.CompleteMultipartUpload(context.Background(), key, v.UploadID, optCom)
	if err != nil {
		return err
	}
	return nil
}
