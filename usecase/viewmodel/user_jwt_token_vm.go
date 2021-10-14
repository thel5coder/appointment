package viewmodel

type UserJwtTokenVm struct {
	Token           string `json:"token"`
	ExpTime         string `json:"exp_time"`
	RefreshToken    string `json:"refresh_token"`
	ExpRefreshToken string `json:"exp_refresh_token"`
	IsActive        bool   `json:"is_active"`
}
