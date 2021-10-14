package viewmodel

type DoctorScheduleVm struct {
	DoctorID   string                 `json:"doctor_id"`
	DoctorName string                 `json:"doctor_name"`
	ClinicID   string                 `json:"clinic_id"`
	ClinicName string                 `json:"clinic_name"`
	Schedules  []DoctorScheduleDaysVm `json:"schedules"`
}

type DoctorScheduleDaysVm struct {
	ID        string                      `json:"id"`
	Day       string                      `json:"day"`
	CreatedAt string                      `json:"created_at"`
	UpdatedAt string                      `json:"updated_at"`
	WorkTimes []DoctorScheduleWorkTimesVm `json:"work_times"`
}

type DoctorScheduleWorkTimesVm struct {
	ID       string `json:"id"`
	DoctorID string `json:"doctor_id"`
	StartAt  string `json:"start_at"`
	EndAt    string `json:"end_at"`
	Status   string `json:"status"`
}

type ScheduleSlotVm struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
