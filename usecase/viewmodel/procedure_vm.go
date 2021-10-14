package viewmodel

type ProcedureVm struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Duration  int    `json:"duration"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
