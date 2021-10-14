package usecase

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"time"
)

type PromotionProductUseCase struct {
	*UcContract
}

//store
func(uc PromotionProductUseCase) Store(productIDs []string, promotionID string) (err error){
	count,err := uc.CountBy("promotion_id",promotionID,"=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-promotionProduct-countByProductID")
		return err
	}

	if count > 0 {
		err = uc.DeleteBy("promotion_id",promotionID,"=")
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-promotionProduct-deleteByProductID")
			return err
		}
	}

	for _, productID := range productIDs{
		err = uc.add(productID,promotionID)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-promotionProduct-add")
			return err
		}
	}

	return nil
}

//add
func (uc PromotionProductUseCase) add(productID, promotionID string) (err error) {
	repository := actions.NewPromotionProductRepository(uc.DB)
	now := time.Now().UTC()

	model := models.PromotionProduct{
		ProductID:   productID,
		PromotionID: promotionID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err = repository.Add(model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-promotionProduct-add")
		return err
	}

	return nil
}

//delete by
func (uc PromotionProductUseCase) DeleteBy(column, value, operator string) (err error) {
	repository := actions.NewPromotionProductRepository(uc.DB)
	now := time.Now().UTC()

	model := models.PromotionProduct{
		UpdatedAt: now,
		DeletedAt: sql.NullTime{Time: now},
	}
	err = repository.DeleteBy(column, value, operator, model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-promotionProduct-deleteBy")
		return err
	}

	return nil
}

//count by
func(uc PromotionProductUseCase) CountBy(column,value,operator string) (res int,err error){
	repository := actions.NewPromotionProductRepository(uc.DB)
	res,err = repository.CountBy(column,value,operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-promotionProduct-countBy")
		return res,err
	}

	return res,nil
}
