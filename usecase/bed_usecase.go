package usecase

import (
	"database/sql"
	"errors"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/helpers/messages"
	"profira-backend/server/requests"
	"profira-backend/usecase/viewmodel"
	"strings"
	"time"
)

type BedUseCase struct {
	*UcContract
}

//browse
func (uc BedUseCase) Browse(search, order, sort string, page, limit int) (res []viewmodel.BedVm, pagination viewmodel.PaginationVm, err error) {
	repository := actions.NewBedRepository(uc.DB)
	offset, limit, page, order, sort := uc.setPaginationParameter(page, limit, order, sort)

	beds, count, err := repository.Browse(uc.ClinicID, search, order, sort, limit, offset)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-browse-bed")
		return res, pagination, err
	}
	for _, bed := range beds {
		res = append(res, uc.buildBody(bed))
	}
	pagination = uc.setPaginationResponse(page, limit, count)

	return res, pagination, nil
}

//browse all
func (uc BedUseCase) BrowseAll(search string) (res []viewmodel.BedVm, err error) {
	repository := actions.NewBedRepository(uc.DB)
	beds, err := repository.BrowseAll(uc.ClinicID, search)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-browseAll-bed")
		return res, err
	}

	for _, bed := range beds {
		res = append(res, uc.buildBody(bed))
	}

	return res, nil
}

//read by
func (uc BedUseCase) ReadBy(column, value, operator string) (res viewmodel.BedVm, err error) {
	repository := actions.NewBedRepository(uc.DB)
	bed, err := repository.ReadBy(column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-readBy-bed")
		return res, err
	}
	res = uc.buildBody(bed)

	return res, nil
}

//edit
func (uc BedUseCase) Edit(input *requests.BedRequest, ID string) (err error) {
	repository := actions.NewBedRepository(uc.DB)

	count, err := uc.countBy(ID, "bed_code", input.BedCode, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-countBy")
		return err
	}
	if count > 0 {
		logruslogger.Log(logruslogger.WarnLevel, messages.DataAlreadyExist, functioncaller.PrintFuncName(), "duplicate-edit-bed")
		return errors.New(messages.DataAlreadyExist)
	}

	now := time.Now().UTC()
	body := models.Bed{
		ID:        ID,
		BedCode:   input.BedCode,
		ClinicID:  uc.ClinicID,
		IsUseAble: input.IsUseAble,
		UpdatedAt: now,
	}
	err = repository.Edit(body,uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-edit-bed")
		return err
	}

	bedTreatmentUc := BedTreatmentUseCase{UcContract:uc.UcContract}
	err = bedTreatmentUc.Store(ID,input.Treatments)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-bedTreatment-store")
		return err
	}

	return nil
}

//add
func (uc BedUseCase) Add(input *requests.BedRequest) (err error) {
	repository := actions.NewBedRepository(uc.DB)

	count, err := uc.countBy("", "bed_code", input.BedCode, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-countBy")
		return err
	}
	if count > 0 {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "duplicate-add-bed")
		return err
	}

	now := time.Now().UTC()
	body := models.Bed{
		BedCode:   input.BedCode,
		ClinicID:  uc.ClinicID,
		IsUseAble: input.IsUseAble,
		CreatedAt: now,
		UpdatedAt: now,
	}
	body.ID, err = repository.Add(body,uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-add-bed")
		return err
	}

	bedTreatmentUc := BedTreatmentUseCase{UcContract:uc.UcContract}
	err = bedTreatmentUc.Store(body.ID	,input.Treatments)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-bedTreatment-store")
		return err
	}

	return nil
}

//delete by
func (uc BedUseCase) DeleteBy(column, value, operator string) (err error) {
	repository := actions.NewBedRepository(uc.DB)

	count, err := uc.countBy("", column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-countBy")
		return err
	}
	if count > 0 {
		now := time.Now().UTC()
		body := models.Bed{
			UpdatedAt: now,
			DeletedAt: sql.NullTime{
				Time: now,
			},
		}
		err = repository.DeleteBy(column, value, operator, body,uc.TX)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-deleteBy")
			return err
		}
	}

	return nil
}

//delete by pk
func(uc BedUseCase) DeleteByPk(ID string) (err error){
	err = uc.DeleteBy("id",ID,"=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-bed-deleteByID")
		return err
	}

	bedTreatmentUc := BedTreatmentUseCase{UcContract:uc.UcContract}
	err = bedTreatmentUc.DeleteBy("bed_id",ID,"=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-bedTreatment-store")
		return err
	}

	return nil
}

func (uc BedUseCase) countBy(ID, column, value, operator string) (res int, err error) {
	repository := actions.NewBedRepository(uc.DB)
	res, err = repository.CountBy(ID, column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-countBy-bed")
		return res, err
	}

	return res, nil
}

//build body
func (uc BedUseCase) buildBody(model models.Bed) viewmodel.BedVm {
	var bedTreatmentVm []viewmodel.BedTreatment
	if model.Treatments.String != "" {
		treatments := strings.Split(model.Treatments.String, ",")
		for _, treatment := range treatments {
			treatmentArr := strings.Split(treatment, ":")
			bedTreatmentVm = append(bedTreatmentVm, viewmodel.BedTreatment{
				ID:          treatmentArr[0],
				TreatmentID: treatmentArr[1],
				Name:        treatmentArr[2],
			})
		}
	}

	return viewmodel.BedVm{
		ID:                model.ID,
		BedCode:           model.BedCode,
		ClinicID:          model.Clinic.ID,
		ClinicName:        model.Clinic.Name,
		ClinicAddress:     model.Clinic.Address,
		ClinicPicName:     model.Clinic.PICName,
		ClinicPhoneNumber: model.Clinic.PhoneNumber,
		ClinicEmail:       model.Clinic.Email,
		IsUseAble:         model.IsUseAble,
		Treatments:        bedTreatmentVm,
		CreatedAt:         model.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         model.UpdatedAt.Format(time.RFC3339),
	}
}
