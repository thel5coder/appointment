package requests

type BedTreatmentRequest struct {
	Selected []string `json:"selected"`
	Deleted  []string `json:"deleted"`
}
