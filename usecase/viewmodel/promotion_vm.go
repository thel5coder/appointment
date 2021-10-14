package viewmodel

import (
	"profira-backend/db/models"
	"time"
)

type PromotionVm struct {
	ID                         string               `json:"id"`
	Slug                       string               `json:"slug"`
	Name                       string               `json:"name"`
	CustomerPromotionCondition string               `json:"customer_promotion_condition"`
	PromotionType              string               `json:"promotion_type"`
	Description                string               `json:"description"`
	StartDate                  string               `json:"start_date"`
	EndDate                    string               `json:"end_date"`
	PhotoPath                  string               `json:"photo_path"`
	NominalType                string               `json:"nominal_type"`
	NominalPercentage          int32                `json:"nominal_percentage"`
	NominalAmount              int32                `json:"nominal_amount"`
	BirthDateCondition         string               `json:"birth_date_condition"`
	SexCondition               string               `json:"sex_condition"`
	RegisterDateConditionStart string               `json:"register_date_condition_start"`
	RegisterDateConditionEnd   string               `json:"register_date_condition_end"`
	Products                   []PromotionProductVm `json:"products"`
	CreatedAt                  string               `json:"created_at"`
	UpdatedAt                  string               `json:"updated_at"`
}

func NewPromotionVm() PromotionVm {
	return PromotionVm{}
}

func (vm PromotionVm) Build(model *models.Promotion) PromotionVm {
	promotionProductVm := NewPromotionProductVm()
	var promotionProduct []PromotionProductVm
	if model.Treatments.String != "" {
		promotionProduct = promotionProductVm.Build(model.Treatments.String)
	}

	return PromotionVm{
		ID:                         model.ID,
		Slug:                       model.Slug,
		Name:                       model.Name,
		CustomerPromotionCondition: model.CustomerPromotionCondition.String,
		PromotionType:              model.PromotionType,
		Description:                model.Description,
		StartDate:                  model.StartDate.Format("2006-01-02"),
		EndDate:                    model.EndDate.Format("2006-01-02"),
		PhotoPath:                  model.FilePath.String,
		NominalType:                model.NominalType,
		NominalPercentage:          model.NominalPercentage.Int32,
		NominalAmount:              model.NominalAmount.Int32,
		BirthDateCondition:         model.BirthDateCondition.Time.Format("2006-01-02"),
		SexCondition:               model.SexCondition.String,
		RegisterDateConditionStart: model.RegisterDateConditionStart.Time.Format(time.RFC3339),
		RegisterDateConditionEnd:   model.RegisterDateConditionEnd.Time.Format(time.RFC3339),
		Products:                   promotionProduct,
		CreatedAt:                  model.CreatedAt.Format(time.RFC3339),
		UpdatedAt:                  model.UpdatedAt.Format(time.RFC3339),
	}
}
