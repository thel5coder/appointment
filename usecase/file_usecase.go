package usecase

import (
	"github.com/globalsign/mgo/bson"
	"mime/multipart"
	"path/filepath"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/usecase/viewmodel"
	"time"
)

type FileUseCase struct {
	*UcContract
}

//read
func (uc FileUseCase) Read(ID string) (res viewmodel.FileVm, err error) {
	repository := actions.NewFileRepository(uc.DB)
	file, err := repository.ReadBy("id", ID, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-file-readBy")
		return res, err
	}
	res = uc.buildBody(file)

	return res, nil
}

//add
func (uc FileUseCase) Add(inputFile *multipart.FileHeader) (res viewmodel.FileVm, err error) {
	repository := actions.NewFileRepository(uc.DB)
	now := time.Now().UTC()

	//upload minio
	fileName := bson.NewObjectId().Hex() + filepath.Ext(inputFile.Filename)
	path := fileName
	minioUc := MinioUseCase{UcContract: uc.UcContract}
	res.Path, err = minioUc.Add(inputFile, path)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-minio-uploadFile")
		return res, err
	}

	//add table file
	model := models.File{
		Path:      path,
		CreatedAt: now,
		UpdatedAt: now,
	}
	res.ID, err = repository.Add(model)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-file-add")
		return res, err
	}

	//get file
	res.Path,err = minioUc.GetFile(res.Path)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-minio-getFile")
		return res, err
	}

	return res, nil
}

//build body
func (uc FileUseCase) buildBody(model models.File) viewmodel.FileVm {
	minioUc := MinioUseCase{UcContract: uc.UcContract}
	path, err := minioUc.GetFile(model.Path)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-minio-getFile")
	}

	return viewmodel.FileVm{
		ID:   model.ID,
		Path: path,
	}
}
