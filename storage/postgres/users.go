package postgres

import (
	"basa/structs"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type usersRepo struct {
	DB *sql.DB
}

func NewUsersRepo(db *sql.DB) usersRepo {
	return usersRepo{
		DB: db,
	}
}

// INSERT USER

func (u usersRepo) InsertUser(user structs.Users) (string, error) {
	fmt.Println("user: ", user)
	id := uuid.New()
	if _, err := u.DB.Exec(`insert into users values ($1,$2,$3,$4,$5)`, id, user.First_Name, user.Last_Name, user.Email, user.Phone); err != nil {
		return "", err
	}
	return id.String(), nil
}

// GET USER BY ID

func (u usersRepo) GetByIDUser(id uuid.UUID) (structs.Users, error) {
	user := structs.Users{}

	if err := u.DB.QueryRow(`select id, first_name, last_name, email, phone from users where id = $1`, id).Scan(
		&user.ID,
		&user.First_Name,
		&user.Last_Name,
		&user.Email,
		&user.Phone,
	); err != nil {
		return structs.Users{}, err
	}
	return user, nil
}

// GET USERS LIST

func (u usersRepo) GetListUser() ([]structs.Users, error) {
	Users := []structs.Users{}

	rows, err := u.DB.Query(`select * from users`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		User := structs.Users{}

		rows.Scan(&User.ID, &User.First_Name, &User.Last_Name, &User.Email, &User.Phone)

		Users = append(Users, User)

	}
	return Users, nil
}

// UPDATE USER

func (u usersRepo) UpdateUser(user structs.Users) {

	_, err := u.DB.Exec(`update users set first_name = $1, last_name = $2, email = $3, phone = $4 where id = $5`, user.First_Name, user.Last_Name, user.Email, user.Phone, user.ID)
	if err != nil {
		log.Fatalln(" error while update car by id: ", err.Error())
	}
	return
}

// DELETE USER

func (u usersRepo) DeleteUser(id uuid.UUID) {
	_, err := u.DB.Exec(`delete from users where id = $1`, id)
	if err != nil {
		log.Fatalln("error while delete user by id")
	}
}
