package repository

import (
	"database/sql"
	"fmt"
	"user-management/internal/errors"
	"user-management/internal/model"
)

type UserRepo interface {
	GetAll() ([]model.User, error)
	GetByID(id int) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Create(user *model.User) error
	Update(id int, user model.User) error
	Delete(id int) error
	ExistsByEmail(email string) bool
	ExistsByID(id int)bool
	ExistsByUsername(username string )bool
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connectionString string) (UserRepo, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to open db conection:%w", err)
	}
	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) GetAll() ([]model.User, error) {
	query := `select id,username,email,name,isactive from Users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error while querying:%w", err)

	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.IsActive)
		if err != nil {
			return nil, fmt.Errorf("error while scaning row data:%w ", err)
		}
		users = append(users, user)
	}
	return users, nil

}

func (r *PostgresRepository) GetByID(id int) (*model.User, error) {
	query := `select id,username,email,password,name,isactive from Users where id=$1`
	row := r.db.QueryRow(query, id)
	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Email,&user.Password, &user.Name, &user.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError(id, "no user with that id")
		}
		return nil, fmt.Errorf("error while scaning row data:%w ", err)
	}
	return &user, nil
}
func (r *PostgresRepository) GetByEmail(email string) (*model.User, error) {
	if r.ExistsByEmail(email){
	query := `select id,username,email,name,isactive,password from Users where email=$1`
	row := r.db.QueryRow(query,email)
	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.IsActive,&user.Password)
	if err != nil {
		return nil, fmt.Errorf("error while scaning row data:%w ", err)
	}
	return &user, nil
}
	return nil, errors.NewNotFoundError(email, "no user with that id")

}
func (r *PostgresRepository) Create(user *model.User) error {
	query := `insert into Users (username,email,password,name) values($1,$2,$3,$4) returning id`
	err := r.db.QueryRow(query, user.Username, user.Email, user.Password,user.Name).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("error while scaning row data:%w ", err)
	}
	return nil
}
func (r *PostgresRepository) Update(id int, user model.User) error {
	//can optimize by calling existsbyid here clear redundant code.
	if r.ExistsByID(id){
	query := `update Users set username=$1 ,email=$2,password=$3,updated_at=CURRENT_TIMESTAMP  where id=$4 `
	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, id)
	if err != nil {
		return fmt.Errorf("unable to exec query %w",err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affectted")
	}
	if rowsAffected == 0 {
		return errors.NewNotFoundError(id, "no user with the given id")
	}
	return nil
}
	return errors.NewNotFoundError(id,"user not found")
}
func (r *PostgresRepository) Delete(id int) error {
	//can optimize by calling existsbyid here clear redundant code.
	query := `delete from Users where id=$1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("unable to exec query")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affectted")
	}
	if rowsAffected == 0 {
		return errors.NewNotFoundError(id, "no user with the given id")
	}
	return nil
}

func (r *PostgresRepository) ExistsByEmail(email string) bool { // Returns bool, not error
	query := `SELECT EXISTS(SELECT 1 FROM Users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false // On error, assume doesn't exist
	}

	return exists
}
func (r *PostgresRepository) ExistsByID(id int) bool { // Returns bool, not error
	query := `SELECT EXISTS(SELECT 1 FROM Users WHERE id = $1)`

	var exists bool
	err := r.db.QueryRow(query,id).Scan(&exists)
	if err != nil {
		return false // On error, assume doesn't exist
	}

	return exists
}
func (r *PostgresRepository) ExistsByUsername(username string) bool {
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
    var exists bool
    err := r.db.QueryRow(query, username).Scan(&exists)
    return err == nil && exists
}
