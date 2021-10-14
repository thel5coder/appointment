package requests

type TreatmentProcedureRequest struct {
	ID          string `json:"id"`
	ProcedureID string `json:"procedure_id" validate:"required"`
	Duration    int32  `json:"duration"`
}
