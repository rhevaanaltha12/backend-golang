package model

type Seller struct {
	Id                uint   `json:"id" gorm:"primary_key"`
	SellerName        string `json:"seller_name" gorm:"size:255;null;"`
	SellerAddress     string `json:"seller_address" gorm:"size:255;null;"`
	SellerPhoneNumber string `json:"seller_phone" gorm:"size:255;null;"`
	SelerDescription  string `json:"seller_description" gorm:"null;"`
	Fish              []Fish `gorm:"many2many:temporary;"`
	ImagePath         string `json:"image_path" gorm:"size:255;null;"`
	CreatedBase
	ArchievedBase
	DeletedBase
}
