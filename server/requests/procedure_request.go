package requests

type ProcedureRequest struct {
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}
