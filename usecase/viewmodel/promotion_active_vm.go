package viewmodel

import "profira-backend/db/models"

type PromotionActiveVm struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	NominalType       string `json:"nominal_type"`
	NominalPercentage int32  `json:"nominal_percentage"`
	NominalAmount     int32  `json:"nominal_amount"`
	FilePath          string `json:"file_path"`
}

func NewPromotionActiveVm() PromotionActiveVm {
	return PromotionActiveVm{}
}

func (vm PromotionActiveVm) Build(model *models.Promotion) PromotionActiveVm {
	return PromotionActiveVm{
		ID:                model.ID,
		Name:              model.Name,
		Description:       model.Description,
		NominalType:       model.NominalType,
		NominalPercentage: model.NominalPercentage.Int32,
		NominalAmount:     model.NominalAmount.Int32,
		FilePath:          model.FilePath.String,
	}
}
