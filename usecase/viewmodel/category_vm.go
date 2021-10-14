package viewmodel

type CategoryVm struct {
	ID             string `json:"id"`
	Slug           string `json:"slug"`
	Name           string `json:"name"`
	ParentID       string `json:"parent_id"`
	IsActive       bool   `json:"is_active"`
	IconPath       string `json:"icon_path"`
	BackgroundPath string `json:"background_path"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type CategoryTreeViewVm struct {
	ID             string                `json:"id"`
	Slug           string                `json:"slug"`
	Name           string                `json:"name"`
	ParentID       string                `json:"parent_id"`
	Treatments     []TreatmentTreeViewVm `json:"treatments"`
	Child          []CategoryTreeViewVm  `json:"child"`
}

type TreatmentTreeViewVm struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
