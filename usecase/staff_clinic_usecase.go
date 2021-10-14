package usecase

import (
	"profira-backend/db/repositories/actions"
	"time"
)

type StaffClinicUseCase struct {
	*UcContract
}

func (uc StaffClinicUseCase) Add(staffID, clinicID string) (err error) {
	repository := actions.NewStaffClinicRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)

	err = repository.Add(staffID, clinicID, now, uc.TX)
	if err != nil {
		return err
	}

	return nil
}

func (uc StaffClinicUseCase) Delete(column, value, operator string) (err error) {
	repository := actions.NewStaffClinicRepository(uc.DB)

	err = repository.DeleteBy(column, value, operator, uc.TX)
	if err != nil {
		return err
	}

	return nil
}

func (uc StaffClinicUseCase) CountBy(column, value, operator string) (res int, err error) {
	repository := actions.NewStaffClinicRepository(uc.DB)
	res, err = repository.CountBy(column, value, operator)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc StaffClinicUseCase) Store(staffID string, clinicIDs []string) (err error) {
	count, err := uc.CountBy("staff_id", staffID, "=")
	if err != nil {
		return err
	}

	if count > 0 {
		err = uc.Delete("staff_id", staffID, "=")
		if err != nil {
			return err
		}
	}

	for _, clinicID := range clinicIDs {
		err = uc.Add(staffID, clinicID)
		if err != nil {
			return err
		}
	}

	return nil
}
