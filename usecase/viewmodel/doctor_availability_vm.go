package viewmodel

//DoctorAvailabilityVm view model section
type DoctorAvailabilityVm struct {
	Id             string                      `json:"id"`
	Name           string                      `json:"name"`
	AvailableSlots []DoctorAvailabilitySlotsVm `json:"available_slots"`
}

type DoctorAvailabilitySlotsVm struct {
	Time      string `json:"time"`
	Available bool   `json:"available"`
}

func NewDoctorAvailabilityVm(id, name string) DoctorAvailabilityVm {
	return DoctorAvailabilityVm{
		Id:             id,
		Name:           name,
		AvailableSlots: nil,
	}
}

//IDoctorAvailability builder section
type IDoctorAvailability interface {
	SetID(id string) IDoctorAvailability

	SetName(name string) IDoctorAvailability

	Build() DoctorAvailabilityVm
}

type DoctorAvailabilityBuilderVm struct {
	id             string
	name           string
	availableSlots []DoctorAvailabilityBuilderSlots
}

type DoctorAvailabilityBuilderSlots struct {
	time      string
	available bool
}

func NewDoctorAvailabilityBuilderVm() IDoctorAvailability{
	return &DoctorAvailabilityBuilderVm{}
}

func (builder *DoctorAvailabilityBuilderVm) SetID(id string) IDoctorAvailability {
	builder.id = id

	return builder
}

func (builder *DoctorAvailabilityBuilderVm) SetName(name string) IDoctorAvailability {
	builder.name = name

	return builder
}

func (builder *DoctorAvailabilityBuilderVm) Build() DoctorAvailabilityVm {
	return NewDoctorAvailabilityVm(builder.id, builder.name)
}
