package requests

type DoctorScheduleRequest struct {
	DoctorID             string                     `json:"doctor_id"`
	ClinicID             string                     `json:"clinic_id"`
	Schedules            []DoctorScheduleDayRequest `json:"schedules"`
	DeletedScheduleDay   []string                   `json:"deleted_schedule_day"`
	DeletedScheduleTimes []string                   `json:"deleted_schedule_times"`
}

type DoctorScheduleDayRequest struct {
	ID            string                      `json:"id"`
	Day           string                      `json:"day"`
	ScheduleTimes []DoctorScheduleTimeRequest `json:"schedule_times"`
}

type DoctorScheduleTimeRequest struct {
	ID      string `json:"id"`
	StartAt string `json:"start_at"`
	EndAt   string `json:"end_at"`
}
