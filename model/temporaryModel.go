package model

type Temporary struct {
	SellerId uint `json:"seller_id" gorm:"null;"`
	FishId   uint `json:"fish_id" gorm:"null;"`
}
