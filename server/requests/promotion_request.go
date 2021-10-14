package requests

type PromotionRequest struct {
	Name                       string   `json:"name" validate:"required"`
	StartDate                  string   `json:"start_date"`
	EndDate                    string   `json:"end_date"`
	CustomerPromotionCondition string   `json:"customer_promotion_condition"`
	Description                string   `json:"description"`
	PromotionType              string   `json:"promotion_type"`
	PhotoID                    string   `json:"photo_id"`
	NominalType                string   `json:"nominal_type"`
	NominalPercentage          int32    `json:"nominal_percentage"`
	NominalAmount              int32    `json:"nominal_amount"`
	BirthDateCondition         string   `json:"birth_date_condition"`
	SexCondition               string   `json:"sex_condition"`
	RegisterDateConditionStart string   `json:"register_date_condition_start"`
	RegisterDateConditionEnd   string   `json:"register_date_condition_end"`
	Treatments                 []string `json:"treatments"`
}
