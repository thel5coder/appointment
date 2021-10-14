package viewmodel

import (
	"profira-backend/helpers/str"
	"strings"
)

type PromotionProductVm struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewPromotionProductVm() PromotionProductVm{
	return PromotionProductVm{}
}

func(vm PromotionProductVm) Build(product string) (res []PromotionProductVm){
	products := str.Unique(strings.Split(product,","))
	for _,product := range products{
		productArr := strings.Split(product,":")
		res = append(res,PromotionProductVm{
			ID:   productArr[0],
			Name: productArr[1],
		})
	}

	return res
}
