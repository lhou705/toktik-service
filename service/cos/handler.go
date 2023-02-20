package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"
	service "toktik/service/cos/kitex_gen/cos"
	"toktik/service/cos/utils"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// CosImpl implements the last service interface defined in the IDL.
type CosImpl struct{}

// Upload implements the CosImpl interface.
func (s *CosImpl) Upload(ctx context.Context, req *service.UploadReq) (resp *service.UploadResp, err error) {
	key := req.GetKey()
	// 构造虚拟文件
	klog.CtxInfof(ctx, "获取到文件%s，准备上传", key)
	file := &utils.File{
		Reader: bytes.NewReader(req.GetFile()),
		FileInfo: utils.FileInfo{
			Data: req.GetFile(),
		},
	}
	resp = &service.UploadResp{}
	v, _, err := client.Object.InitiateMultipartUpload(ctx, key, nil)
	if err != nil {
		errStr := fmt.Sprintf("初始化上传文件错误，文件名：%s，错误原因：%v", key, err)
		klog.CtxErrorf(ctx, errStr)
		resp.IsSuccess = false
		return
	}
	partNum, partSize := cos.DividePart(file.Size(), 16)
	klog.CtxInfof(ctx, "文件分块成功，总计%d块，每块大小%d", partNum+1, partSize)

	var chunk cos.Chunk
	var chunks []cos.Chunk
	for i := int64(0); i < partNum; i++ {
		chunk.Number = int(i + 1)
		chunk.OffSet = i * partSize
		chunk.Size = partSize
		chunks = append(chunks, chunk)
	}

	if file.Size()%partSize > 0 {
		chunk.Number = len(chunks) + 1
		chunk.OffSet = int64(len(chunks)) * partSize
		chunk.Size = file.Size() % partSize
		chunks = append(chunks, chunk)
		partNum++
	}
	optCom := &cos.CompleteMultipartUploadOptions{}
	var wg sync.WaitGroup
	for _, chunk := range chunks {
		chunk := chunk
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := file.Seek(chunk.OffSet, io.SeekStart)
			if err != nil {
				errStr := fmt.Sprintf("读取文件块%d失败，原因：%v", chunk.Number, err)
				klog.CtxErrorf(ctx, errStr)
				resp.IsSuccess = false
				return
			}
			uploadResp, err := client.Object.UploadPart(
				context.Background(),
				key,
				v.UploadID,
				chunk.Number,
				cos.LimitReadCloser(file, partSize),
				&cos.ObjectUploadPartOptions{ContentLength: chunk.Size})
			if err != nil {
				errStr := fmt.Sprintf("上传文件块%d失败，原因：%v", chunk.Number, err)
				klog.CtxErrorf(ctx, errStr)
				resp.IsSuccess = false
				return
			}
			optCom.Parts = append(optCom.Parts, cos.Object{ETag: uploadResp.Header.Get("ETag"), PartNumber: chunk.Number})
			time.Sleep(time.Millisecond)
		}()
	}
	wg.Wait()
	_, _, err = client.Object.CompleteMultipartUpload(context.Background(), key, v.UploadID, optCom)
	if err != nil {
		errStr := fmt.Sprintf("上传文件%s失败，原因：%v", key, err)
		klog.CtxErrorf(ctx, errStr)
		resp.IsSuccess = false
		return
	}
	klog.CtxInfof(ctx, "上传成功")
	resp.IsSuccess = true
	return
}
