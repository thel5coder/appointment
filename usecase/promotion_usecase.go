package usecase

import (
	"database/sql"
	"github.com/gosimple/slug"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/datetime"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/server/requests"
	"profira-backend/usecase/viewmodel"
	"time"
)

type PromotionUseCase struct {
	*UcContract
}

//browse
func (uc PromotionUseCase) Browse(search, order, sort string, page, limit int) (res []viewmodel.PromotionVm, pagination viewmodel.PaginationVm, err error) {
	repository := actions.NewPromotionRepository(uc.DB)
	offset, limit, page, order, sort := uc.setPaginationParameter(page, limit, order, sort)

	promotions, count, err := repository.Browse(search, order, sort, limit, offset)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-promotion-browse")
		return res, pagination, err
	}

	minioUc := MinioUseCase{UcContract: uc.UcContract}
	vm := viewmodel.NewPromotionVm()
	for _, promotion := range promotions {
		if promotion.FilePath.String != "" {
			promotion.FilePath.String, err = minioUc.GetFile(promotion.FilePath.String)
			if err != nil {
				logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-minio-getFile")
				return res, pagination, err
			}
		}
		res = append(res, vm.Build(&promotion))
	}
	pagination = uc.setPaginationResponse(page, limit, count)

	return res, pagination, err
}

//browse active promotion
func (uc PromotionUseCase) BrowseActivePromotion() (res []viewmodel.PromotionActiveVm, err error) {
	repository := actions.NewPromotionRepository(uc.DB)
	vm := viewmodel.NewPromotionActiveVm()
	now := time.Now().UTC().Format("2006-01-02")
	date := datetime.StrParseToTime(now,"2006-01-02")

	promotions, err := repository.BrowseActivePromotion(date.Unix(), date.Unix())
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-promotion-browseActivePromotion")
		return res, err
	}

	minioUc := MinioUseCase{UcContract: uc.UcContract}
	for _, promotion := range promotions {
		promotion.FilePath.String, err = minioUc.GetFile(promotion.FilePath.String)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-minio-getFile")
		}
		res = append(res, vm.Build(&promotion))
	}

	return res, nil
}

//read
func (uc PromotionUseCase) ReadBy(column, value, operator string) (res viewmodel.PromotionVm, err error) {
	repository := actions.NewPromotionRepository(uc.DB)
	promotion, err := repository.ReadBy(column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-promotion-readBy")
		return res, err
	}

	minioUc := MinioUseCase{UcContract: uc.UcContract}
	if promotion.FilePath.String != "" {
		promotion.FilePath.String, err = minioUc.GetFile(promotion.FilePath.String)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-minio-getFile")
			return res, err
		}
	}

	vm := viewmodel.NewPromotionVm()
	res = vm.Build(&promotion)

	return res, nil
}

//edit
func (uc PromotionUseCase) Edit(input *requests.PromotionRequest, ID string) (err error) {
	repository := actions.NewPromotionRepository(uc.DB)
	now := time.Now().UTC()

	model := models.Promotion{
		ID:                         ID,
		Slug:                       slug.Make(input.Name),
		Name:                       input.Name,
		CustomerPromotionCondition: sql.NullString{String: input.CustomerPromotionCondition},
		PromotionType:              input.PromotionType,
		Description:                input.Description,
		StartDate:                  datetime.StrParseToTime(input.StartDate, "2006-01-02"),
		StartAtUnix:                datetime.StrParseToTime(input.StartDate, "2006-01-02").Unix(),
		EndAtUnix:                  datetime.StrParseToTime(input.EndDate, "2006-01-02").Unix(),
		EndDate:                    datetime.StrParseToTime(input.EndDate, "2006-01-02"),
		FotoID:                     sql.NullString{String: input.PhotoID},
		NominalType:                input.NominalType,
		NominalPercentage:          sql.NullInt32{Int32: input.NominalPercentage},
		NominalAmount:              sql.NullInt32{Int32: input.NominalAmount},
		BirthDateCondition:         sql.NullTime{Time: datetime.StrParseToTime(input.BirthDateCondition, "2006-01-02")},
		SexCondition:               sql.NullString{String: input.SexCondition},
		RegisterDateConditionStart: sql.NullTime{Time: datetime.StrParseToTime(input.RegisterDateConditionStart, "2006-01-02")},
		RegisterDateConditionEnd:   sql.NullTime{Time: datetime.StrParseToTime(input.RegisterDateConditionEnd, "2006-01-02")},
		UpdatedAt:                  now,
	}
	err = repository.Edit(model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-promotion-edit")
		return err
	}

	promotionProductUc := PromotionProductUseCase{UcContract: uc.UcContract}
	err = promotionProductUc.Store(input.Treatments, ID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-promotionProduct-store")
		return err
	}

	return nil
}

//add
func (uc PromotionUseCase) Add(input *requests.PromotionRequest) (err error) {
	repository := actions.NewPromotionRepository(uc.DB)
	now := time.Now().UTC()

	model := models.Promotion{
		Slug:                       slug.Make(input.Name),
		Name:                       input.Name,
		CustomerPromotionCondition: sql.NullString{String: input.CustomerPromotionCondition},
		PromotionType:              input.PromotionType,
		Description:                input.Description,
		StartDate:                  datetime.StrParseToTime(input.StartDate, "2006-01-02"),
		StartAtUnix:                datetime.StrParseToTime(input.StartDate, "2006-01-02").Unix(),
		EndAtUnix:                  datetime.StrParseToTime(input.EndDate, "2006-01-02").Unix(),
		EndDate:                    datetime.StrParseToTime(input.EndDate, "2006-01-02"),
		FotoID:                     sql.NullString{String: input.PhotoID},
		NominalType:                input.NominalType,
		NominalPercentage:          sql.NullInt32{Int32: input.NominalPercentage},
		NominalAmount:              sql.NullInt32{Int32: input.NominalAmount},
		BirthDateCondition:         sql.NullTime{Time: datetime.StrParseToTime(input.BirthDateCondition, "2006-01-02")},
		SexCondition:               sql.NullString{String: input.SexCondition},
		RegisterDateConditionStart: sql.NullTime{Time: datetime.StrParseToTime(input.RegisterDateConditionStart, "2006-01-02")},
		RegisterDateConditionEnd:   sql.NullTime{Time: datetime.StrParseToTime(input.RegisterDateConditionEnd, "2006-01-02")},
		CreatedAt:                  now,
		UpdatedAt:                  now,
	}
	model.ID,err = repository.Add(model,uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-promotion-add")
		return err
	}

	promotionProductUc := PromotionProductUseCase{UcContract:uc.UcContract}
	err = promotionProductUc.Store(input.Treatments,model.ID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-promotionProduct-store")
		return err
	}

	return nil
}

//delete
func(uc PromotionUseCase) Delete(ID string) (err error){
	repository := actions.NewPromotionRepository(uc.DB)
	now := time.Now().UTC()

	model := models.Promotion{
		ID: ID,
		UpdatedAt: now,
		DeletedAt: sql.NullTime{Time: now},
	}
	err = repository.Delete(model,uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-promotion-delete")
		return err
	}

	promotionProductUc := PromotionProductUseCase{UcContract:uc.UcContract}
	err = promotionProductUc.DeleteBy("promotion_id",ID,"=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-promotionProduct-deleteByPromotionProductID")
		return err
	}

	return nil
}
