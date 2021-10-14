package viewmodel

type BedVm struct {
	ID                string         `json:"id"`
	BedCode           string         `json:"bed_code"`
	ClinicID          string         `json:"clinic_id"`
	ClinicName        string         `json:"clinic_name"`
	ClinicAddress     string         `json:"clinic_address"`
	ClinicPicName     string         `json:"clinic_pic_name"`
	ClinicPhoneNumber string         `json:"clinic_phone_number"`
	ClinicEmail       string         `json:"clinic_email"`
	IsUseAble         bool           `json:"is_use_able"`
	Treatments        []BedTreatment `json:"treatments"`
	CreatedAt         string         `json:"created_at"`
	UpdatedAt         string         `json:"updated_at"`
}

type BedTreatment struct {
	ID          string `json:"id"`
	TreatmentID string `json:"treatment_id"`
	Name        string `json:"name"`
}
