package daos

import (
"../models"
)

// UserDao DAO interface definition
type UserDao interface {
	Get(i int) (models.User, error)
	GetAll() ([]models.User, error)
	Create(u *models.User) error
	Delete(i int) error
	Update(u *models.User) error
	CheckDuplicateUser(name string) bool
	
	//DeleteDepartmentByUserId(id int) error
}
