package viewmodel

type DoctorVm struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
	User        UserStaffVm         `json:"user"`
	Clinics     []StaffClinicVm     `json:"clinics"`
	Treatments  []DoctorTreatmentVm `json:"treatments"`
}

type DoctorTreatmentVm struct {
	ID          string `json:"id"`
	TreatmentID string `json:"treatment_id"`
	Name        string `json:"name"`
}
