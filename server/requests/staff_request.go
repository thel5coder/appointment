package requests

type StaffRequest struct {
	Name             string                 `json:"name" validate:"required"`
	MobilePhone      string                 `json:"mobile_phone" validate:"required"`
	Email            string                 `json:"email" validate:"required"`
	Password         string                 `json:"password"`
	Description      string                 `json:"description"`
	IsActive         bool                   `json:"is_active"`
	RoleID           string                 `json:"role_id"`
	ProfilePictureID string                 `json:"profile_picture_id"`
	ClinicIDs        []string               `json:"clinic_ids"`
	Treatment        DoctorTreatmentRequest `json:"treatment"`
}
