package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID            int64  `json:"id"`
	Email         string `json:"email"           gorm:"unique"`
	Username      string `json:"username"        gorm:"unique"`
	Password      string `json:"password"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	PlaceOfLiving string `json:"place_of_living"`
	Role          string `json:"role"`
}
