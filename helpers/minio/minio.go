package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"net/url"
	"time"
)

type IMinio interface {
	Upload(bucketName, path string, file *multipart.FileHeader) (res string, err error)
	GetFile(bucketName, objectName string) (res string, err error)
	Delete(bucketName, objectName string) (err error)
}

type minioModel struct {
	Client *minio.Client
}

const defaultDuration = 15

func NewMinioModel(client *minio.Client) IMinio {
	return &minioModel{Client: client}
}

func (model minioModel) Upload(bucketName,path string, fileHeader *multipart.FileHeader) (res string, err error) {
	src, err := fileHeader.Open()
	if err != nil {
		return res, err
	}
	contentType := fileHeader.Header.Get("Content-Type")
	src.Close()

	_, err = model.Client.PutObject(context.Background(), bucketName, path, src, fileHeader.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return res, err
	}
	res = path

	return res, nil
}

func (model minioModel) GetFile(bucketName, objectName string) (res string, err error) {
	reqParams := make(url.Values)

	duration := time.Minute * defaultDuration
	uri, err := model.Client.PresignedGetObject(context.Background(), bucketName, objectName, duration, reqParams)
	if err != nil {
		return res, err
	}
	res = uri.String()

	return res, err
}

func (model minioModel) Delete(bucketName, objectName string) (err error) {
	options := minio.RemoveObjectOptions{}
	err = model.Client.RemoveObject(context.Background(), bucketName, objectName, options)

	return err
}
