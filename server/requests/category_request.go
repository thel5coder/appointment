package requests

type CategoryRequest struct {
	Name             string `json:"name"`
	ParentID         string `json:"parent_id"`
	IsActive         bool   `json:"is_active"`
	FileIconID       string `json:"file_icon_id"`
	FileBackgroundID string `json:"file_background_id"`
}
