package model

type Fish struct {
	Id              uint     `json:"id" gorm:"primary_key"`
	FishName        string   `json:"fish_name" gorm:"size:255;null;"`
	FishType        string   `json:"fish_type" gorm:"size:255;null;"`
	FishSize        string   `json:"fish_size" gorm:"size:255;null;"`
	FishDescription string   `json:"fish_description" gorm:"null;"`
	FishPrice       int      `json:"fish_price" gorm:"null;"`
	Seller          []Seller `gorm:"many2many:temporary;"`
	ImagePath       string   `json:"image_path" gorm:"size:255;null;"`
	CreatedBase
	ArchievedBase
	DeletedBase
}
