package requests

type TreatmentRequest struct {
	Name              string                      `json:"name"`
	Description       string                      `json:"description"`
	Price             int32                       `json:"price"`
	PhotoID           string                      `json:"photo_id"`
	IconID            string                      `json:"icon_id"`
	Procedures        []TreatmentProcedureRequest `json:"procedures"`
	DeletedProcedures []string                    `json:"deleted_procedures"`
	CategoryID        string                      `json:"category_id"`
}
