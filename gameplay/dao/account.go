package dao

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Account  string `gorm:"uniqueIndex;size:64"`
	Password string `gorm:"size:128"`
}

func (a Account) DBId() uint64 {
	return uint64(a.ID)
}

type RoleInfo struct {
	gorm.Model
	Account   string `gorm:"index;size:64"`
	RoleName  string `gorm:"uniqueIndex;size:64"`
	RoleIndex int    `gorm:"size:2"`
}

func (r RoleInfo) DBId() uint64 {
	return uint64(r.ID)
}
