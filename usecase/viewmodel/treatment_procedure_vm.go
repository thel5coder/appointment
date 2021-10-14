package viewmodel

type TreatmentProcedureVm struct {
	ID                string `json:"id"`
	ProcedureID       string `json:"procedure_id"`
	ProcedureName     string `json:"procedure_name"`
	ProcedureDuration int    `json:"procedure_duration"`
}
