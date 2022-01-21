package pkg

import (
	"context"
	"os"
	"path"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiNiuClient struct {
	AccessKey     string `json:"access_key"`
	SecretKey     string `json:"secret_key"`
	Bucket        string `json:"bucket"`
	UseHTTPS      bool   `json:"use_https"`
	UseCdnDomains bool   `json:"use_cdn_domains"`
	Domain        string `json:"domain"`
	Subdir        string `json:"subdir"`
}

func NewQiNiuClient(accessKey, secretKey, bucket string, useHttps, useCdnDomains bool, domain, subdir string) *QiNiuClient {
	qiNiuClient := &QiNiuClient{
		AccessKey:     accessKey,
		SecretKey:     secretKey,
		Bucket:        bucket,
		UseHTTPS:      useHttps,
		UseCdnDomains: useCdnDomains,
		Domain:        domain,
		Subdir:        subdir,
	}
	if !strings.HasSuffix(qiNiuClient.Domain, "/") {
		qiNiuClient.Domain = qiNiuClient.Domain + "/"
	}
	return qiNiuClient
}

func (q *QiNiuClient) UploadImages(files []string) (urls []string) {

	putPolicy := storage.PutPolicy{Scope: q.Bucket}
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	token := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	cfg.UseHTTPS = q.UseHTTPS
	cfg.UseCdnDomains = q.UseCdnDomains

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	for _, file := range files {
		_, fileName := path.Split(file)
		key := q.Subdir + "/" + time.Now().Format("060102-150405") + "-" + fileName

		logger.Println("Start uploading", file, "as", key)
		f, err := os.Open(file)
		if err != nil {
			return
		}
		defer f.Close()

		data, size, err := CompressImg(f)
		if err != nil {
			logger.Fatalln("Error:", err)
			return
		}

		if size == 0 {
			f.Seek(0, 0)
			info, _ := f.Stat()
			size = info.Size()
		}

		if err := formUploader.Put(context.Background(), &ret, token, key, data, size, nil); err != nil {
			logger.Fatalln("Error:", err)
		}

		urls = append(urls, q.Domain+key)
	}
	return
}
