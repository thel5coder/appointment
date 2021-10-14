package viewmodel

type CustomerVm struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	Sex                string         `json:"sex"`
	Address            string         `json:"address"`
	BirthDate          string         `json:"birth_date"`
	MaritalStatus      string         `json:"marital_status"`
	PhoneNumber        string         `json:"phone_number"`
	MobilePhoneNumber1 string         `json:"mobile_phone_number_1"`
	MobilePhoneNumber2 string         `json:"mobile_phone_number_2"`
	Religion           string         `json:"religion"`
	Education          string         `json:"education"`
	Hobby              string         `json:"hobby"`
	Profession         string         `json:"profession"`
	Reference          string         `json:"reference"`
	Notes              string         `json:"notes"`
	CityID             string         `json:"city_id"`
	CreatedAt          string         `json:"created_at"`
	UpdatedAt          string         `json:"updated_at"`
	DeletedAt          string         `json:"deleted_at"`
	User               UserCustomerVm `json:"user"`
}

type UserCustomerVm struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	MobilePhone    string `json:"mobile_phone"`
	RoleID         string `json:"role_id"`
	ProfilePicture FileVm `json:"profile_picture"`
	IsActive       bool   `json:"is_active"`
}
