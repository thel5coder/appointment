package usecase

import (
	"database/sql"
	"github.com/gosimple/slug"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/helpers/str"
	"profira-backend/server/requests"
	"profira-backend/usecase/viewmodel"
	"strings"
	"time"
)

type CategoryUseCase struct {
	*UcContract
}

//browse tree view
func (uc CategoryUseCase) BrowseTreeView() (res []viewmodel.CategoryTreeViewVm, err error) {
	repository := actions.NewCategoryRepository(uc.DB)

	categories, err := repository.BrowseTree("")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-category-browseTreeView")
		return res, err
	}

	res = uc.buildBodyTree(categories)

	return res, nil
}

//browse child by parent slug
func (uc CategoryUseCase) BrowseByParent(ID string) (res []viewmodel.CategoryVm, err error) {
	repository := actions.NewCategoryRepository(uc.DB)

	categories, err := repository.BrowseByParentID(ID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-category-browseChildByParentSlug")
		return res, err
	}

	for _, category := range categories {
		res = append(res, uc.buildBody(category))
	}

	return res, nil
}

//browse shortcut category
func (uc CategoryUseCase) BrowseShortcutCategory() (res []viewmodel.CategoryVm, err error) {
	medicalCategory, err := uc.ReadBy("c.slug", "medical-treatment", "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-category-readBySlugMedical")
		err = nil
	}

	medicalCategories, err := uc.BrowseByParent(medicalCategory.ID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-category-browseChildByParentMedicalTreatment")
		err = nil
	}

	for _, medicalCategory := range medicalCategories {
		res = append(res, medicalCategory)
	}

	beautyCategory, err := uc.ReadBy("c.slug", "beauty-treatment", "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-category-readBySlugBeauty")
		err = nil
	}

	beautyCategories, err := uc.BrowseByParent(beautyCategory.ID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-category-browseChildByParentBeautyTreatment")
		err = nil
	}
	for _, beautyCategory := range beautyCategories {
		res = append(res, beautyCategory)
	}

	return res, nil
}

//read
func (uc CategoryUseCase) ReadBy(column, value, operator string) (res viewmodel.CategoryVm, err error) {
	repository := actions.NewCategoryRepository(uc.DB)

	category, err := repository.ReadBy(column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-category-readBy")
		return res, err
	}
	res = uc.buildBody(category)

	return res, nil
}

//edit
func (uc CategoryUseCase) Edit(input *requests.CategoryRequest, ID string) (err error) {
	repository := actions.NewCategoryRepository(uc.DB)
	now := time.Now().UTC()

	model := models.Category{
		ID:               ID,
		Slug:             slug.Make(input.Name),
		Name:             input.Name,
		ParentID:         sql.NullString{String: input.ParentID},
		IsActive:         input.IsActive,
		FileIconID:       sql.NullString{String: input.FileIconID},
		FileBackgroundID: sql.NullString{String: input.FileBackgroundID},
		UpdatedAt:        now,
	}
	_, err = repository.Edit(model)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-category-edit")
		return err
	}

	return nil
}

//add
func (uc CategoryUseCase) Add(input *requests.CategoryRequest) (err error) {
	repository := actions.NewCategoryRepository(uc.DB)
	now := time.Now().UTC()

	model := models.Category{
		Slug:             slug.Make(input.Name),
		Name:             input.Name,
		ParentID:         sql.NullString{String: input.ParentID},
		IsActive:         input.IsActive,
		FileIconID:       sql.NullString{String: input.FileIconID},
		FileBackgroundID: sql.NullString{String: input.FileBackgroundID},
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	_, err = repository.Add(model)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-category-edit")
		return err
	}

	return nil
}

//delete
func (uc CategoryUseCase) Delete(ID string) (err error) {
	repository := actions.NewCategoryRepository(uc.DB)
	now := time.Now().UTC()

	model := models.Category{
		ID:        ID,
		UpdatedAt: now,
		DeletedAt: sql.NullTime{Time: now},
	}
	_, err = repository.Delete(model)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-category-delete")
		return err
	}

	return nil
}

//build body
func (uc CategoryUseCase) buildBody(model models.Category) viewmodel.CategoryVm {
	minioUc := MinioUseCase{UcContract: uc.UcContract}
	var err error

	var iconPath string
	if model.FileIconPath.String != "" {
		iconPath, err = minioUc.GetFile(model.FileIconPath.String)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-minio-getFileIconPath")
		}
	}

	var backgroundPath string
	if model.FileBackgroundPath.String != "" {
		backgroundPath, err = minioUc.GetFile(model.FileBackgroundPath.String)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-minio-getFileBackgroundPath")
		}
	}

	return viewmodel.CategoryVm{
		ID:             model.ID,
		Slug:           model.Slug,
		Name:           model.Name,
		ParentID:       model.ParentID.String,
		IsActive:       model.IsActive,
		IconPath:       iconPath,
		BackgroundPath: backgroundPath,
		CreatedAt:      model.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      model.UpdatedAt.Format(time.RFC3339),
	}
}

//build body tree
func (uc CategoryUseCase) buildBodyTree(models []models.Category) (res []viewmodel.CategoryTreeViewVm) {
	for _, model := range models {
		if model.ParentID.String == "" {
			res = append(res, uc.buildParentCategory(model, models))
		}
	}

	return res
}

//build parent
func (uc CategoryUseCase) buildParentCategory(model models.Category, models []models.Category) (res viewmodel.CategoryTreeViewVm) {
	res = viewmodel.CategoryTreeViewVm{
		ID:         model.ID,
		Slug:       model.Slug,
		Name:       model.Name,
		ParentID:   model.ParentID.String,
		Treatments: nil,
		Child:      uc.buildChild(model.ID, models),
	}
	return res
}

//build child
func (uc CategoryUseCase) buildChild(parentID string, models []models.Category) (res []viewmodel.CategoryTreeViewVm) {
	for _, model := range models {
		if parentID == model.ParentID.String {
			res = append(res, viewmodel.CategoryTreeViewVm{
				ID:         model.ID,
				Slug:       model.Slug,
				Name:       model.Name,
				ParentID:   model.ParentID.String,
				Treatments: uc.buildCategoryTreatments(model.Treatments.String),
				Child:      uc.buildChild(model.ID, models),
			})
		}
	}

	return res
}

//build category treatments
func (uc CategoryUseCase) buildCategoryTreatments(treatment string) (res []viewmodel.TreatmentTreeViewVm) {
	if treatment != "" {
		treatments := str.Unique(strings.Split(treatment, ","))
		for _, treatment := range treatments {
			treatmentArr := strings.Split(treatment, ":")
			res = append(res, viewmodel.TreatmentTreeViewVm{
				ID:   treatmentArr[0],
				Name: treatmentArr[1],
			})
		}
	}

	return res
}
