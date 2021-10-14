package requests

type DoctorTreatmentRequest struct {
	Selected []string `json:"selected"`
	Deleted  []string `json:"deleted"`
}
