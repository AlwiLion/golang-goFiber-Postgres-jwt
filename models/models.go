package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" gorm:"unique" validate:"required,email"`
	Password []byte `json:"-"`
	IsKyc    bool   `json:"isKyc"`
	Mobile   string `json:"mobile" validate:"required" gorm:"unique"`
	Pan      string `json:"pan" validate:"required" gorm:"unique"`
}

type UserStatment struct {
	gorm.Model
	UserID    uint    `json:"userID"`
	TxnId     string  `json:"txnID"`
	Amount    float64 `json:"amount"`
	IsDeposit bool    `json:"isDeposit"`
	Comment   string  `json:"comment"`
}
