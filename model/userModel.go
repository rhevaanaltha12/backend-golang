package model

import "time"

type User struct {
	Id               uint      `json:"id" gorm:"primary_key"`
	Fullname         string    `json:"fullname" gorm:"size:255;null" validate:"required"`
	Email            string    `json:"email" gorm:"size:255;null;unique" validate:"required"`
	EmailVerifiedAt  time.Time `json:"email_verified_at" gorm:"null"`
	Password         string    `json:"password" gorm:"size:255;null;" validate:"required"`
	RememberToken    string    `json:"remember_token" gorm:"size:255;null;"`
	PasswordExpires  time.Time `json:"password_expires" gorm:"null;"`
	PhoneNumber      string    `json:"phone_number" gorm:"size:255;null;" validate:"required"`
	Role             string    `json:"role" gorm:"size:255;null;type:enum('Admin','User');default:User;" validate:"required"`
	ChangePasswordAt time.Time `json:"change_password_at" gorm:"null;"`
	CreatedBase
	ArchievedBase
	DeletedBase
}

type Register struct {
	Fullname    string `json:"fullname" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Role        string `json:"role" validate:"required"`
}

type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
