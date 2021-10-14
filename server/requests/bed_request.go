package requests

type BedRequest struct {
	BedCode    string              `json:"bed_code"`
	IsUseAble  bool                `json:"is_use_able"`
	Treatments BedTreatmentRequest `json:"treatments"`
}
