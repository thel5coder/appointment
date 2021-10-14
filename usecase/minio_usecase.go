package usecase

import (
	"mime/multipart"
	"os"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/helpers/minio"
)

type MinioUseCase struct {
	*UcContract
}

//add
func (uc MinioUseCase) Add(file *multipart.FileHeader, path string) (res string, err error) {
	minioModel := minio.NewMinioModel(uc.MinioClient)
	res, err = minioModel.Upload(os.Getenv("MINIO_BUCKET"), path, file)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "minio-uploadFile")
		return res, err
	}

	return res, nil
}

//get file
func (uc MinioUseCase) GetFile(path string) (res string, err error) {
	minioModel := minio.NewMinioModel(uc.MinioClient)
	res, err = minioModel.GetFile(os.Getenv("MINIO_BUCKET"), path)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "minio-getFile")
		return res, err
	}

	return res, err
}
