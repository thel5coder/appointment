package usecase

import "profira-backend/helpers/enums/sexenum"

type EnumUseCase struct {
	*UcContract
}

func (uc EnumUseCase) GetSexEnum() []map[string]interface{}{
	var res []map[string]interface{}
	sexEnums := sexenum.GetEnums()

	for _,sexEnum := range sexEnums{
		res = append(res,map[string]interface{}{
			"key":sexEnum["key"],
			"text":sexEnum["text"],
		})
	}

	return res
}



