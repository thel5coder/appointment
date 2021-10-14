package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)


type Connection struct {
	AccessKey string
	SecretKey string
	UseSSL    string
	EndPoint  string
	Duration  int
}

func (conn Connection) InitClient() (client *minio.Client, err error) {
	useSSl := true
	if conn.UseSSL == "false" {
		useSSl = false
	}

	client, err = minio.New(conn.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conn.AccessKey, conn.SecretKey, ""),
		Secure: useSSl,
	})
	if err != nil {
		return client, err
	}

	return client, nil
}
