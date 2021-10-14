package viewmodel

type TreatmentVm struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Duration    int32                  `json:"duration"`
	Price       int32                  `json:"price"`
	PhotoID     string                 `json:"photo_id"`
	PhotoPath   string                 `json:"photo_path"`
	IconID      string                 `json:"icon_id"`
	IconPath    string                 `json:"icon_path"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	Procedures  []TreatmentProcedureVm `json:"procedures"`
}
