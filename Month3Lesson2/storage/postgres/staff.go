package postgres

import (
	"database/sql"
	"fmt"
	"mini_market/api/models"
	"mini_market/storage"

	"github.com/google/uuid"
)

type staffRepo struct {
	db *sql.DB
}

func NewStaffRepo(db *sql.DB) storage.IStaffRepo {
	return &staffRepo{
		db: db,
	}
}

func (s *staffRepo) Create(staff models.CreateStaff) (string, error) {

	uid := uuid.New()

	if _, err := s.db.Exec(`
	 insert into staff values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	 `,
		uid,
		staff.BranchID,
		staff.TarifID,
		staff.StaffType,
		staff.Name,
		staff.Balance,
		staff.BirthDate,
		staff.Gender,
		staff.Login,
		staff.Password,
	); err != nil {
		fmt.Println("error while inserting data", err.Error())
		return "", err
	}

	return "", nil
}

func (s *staffRepo) Get(id string) (models.Staff, error) {

	staff := models.Staff{}

	query := `
		select id, branch_id, tarif_id, staff_type, name, balance, birth_date, gender, login, password, created_at, updated_at, deleted_at from staff where id = $1
`
	if err := s.db.QueryRow(query, id).Scan(
		&staff.ID,
		&staff.BranchID,
		&staff.TarifID,
		&staff.StaffType,
		&staff.Name,
		&staff.Balance,
		&staff.BirthDate,
		&staff.Gender,
		&staff.Login,
		&staff.Password,
		&staff.CreatedAt,
		&staff.UpdatedAt,
		&staff.DeletedAt,
	); err != nil {
		fmt.Println("error while scanning user", err.Error())
		return models.Staff{}, err
	}

	return staff, nil
}

func (s *staffRepo) GetList(req models.GetListRequest) (models.StaffResponse, error) {

	var (
		staffs            = []models.Staff{}
		count             = 0
		countQuery, query string
		page              = req.Page
		offset            = (page - 1) * req.Limit
	)

	countQuery = `
	SELECT count(1) from staff `

	if err := s.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of users", err.Error())
		return models.StaffResponse{}, err
	}

	query = `
	SELECT  id, branch_id, tarif_id, staff_type, name, balance, birth_date, gender, login, password, created_at, updated_at, deleted_at 
		FROM staff
			`

	query += ` LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.StaffResponse{}, err
	}

	for rows.Next() {
		staff := models.Staff{}

		if err = rows.Scan(
			&staff.ID,
			&staff.BranchID,
			&staff.TarifID,
			&staff.StaffType,
			&staff.Name,
			&staff.Balance,
			&staff.BirthDate,
			&staff.Gender,
			&staff.Login,
			&staff.Password,
			&staff.CreatedAt,
			&staff.UpdatedAt,
			&staff.DeletedAt,
		); err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.StaffResponse{}, err
		}

		staffs = append(staffs, staff)
	}

	return models.StaffResponse{
		Staff: staffs,
		Count: count,
	}, nil

}

func (s *staffRepo) Update(staff models.UpdateStaff) (string, error) {
	query := `
	update staff 
		set  branch_id = $1, tarif_id = $2, staff_type = $3, name = $4, balance = $5, birth_date = $6, gender = $7, login = $8, password = $9
			where id = $10`

	if _, err := s.db.Exec(query, staff.BranchID, staff.TarifID, staff.StaffType, staff.Name, staff.Balance, staff.BirthDate, staff.Gender, staff.Login, staff.Password, staff.ID); err != nil {
		fmt.Println("error while updating staff data", err.Error())
		return "", err
	}

	return staff.ID, nil
}

func (s *staffRepo) Delete(delete models.DeleteStaff) error {

	query := `
	delete from staff
		where id = $1
`
	if _, err := s.db.Exec(query, delete.ID); err != nil {
		fmt.Println("error while deleting customer by id", err.Error())
		return err
	}

	return nil
}
