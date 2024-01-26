package storage

import "mini_market/api/models"

type IStorage interface {
	CloseDB()
	Staff() IStaffRepo
}

type IStaffRepo interface {
	Create(models.CreateStaff) (string, error)
	Get(string) (models.Staff, error)
	GetList(models.GetListRequest) (models.StaffResponse, error)
	Update(models.UpdateStaff) (string, error)
	Delete(models.DeleteStaff) error
}
