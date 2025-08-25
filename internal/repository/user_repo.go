package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
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
	ExistsByID(id int) bool
	ExistsByUsername(username string) bool
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connectionString string) (UserRepo, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		slog.Error("unable to open db conection", "error", err, "driver", "postgres")
		return nil, err
	}
	if err := db.Ping(); err != nil {
		slog.Error("database ping failed", "error", err)
		return nil, err
	}
	slog.Info("database conn established - ping success!")
	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) GetAll() ([]model.User, error) {
	query := `select id,username,email,name,isactive from Users`
	rows, err := r.db.Query(query)
	if err != nil {
		slog.Error("failed to execute GetAll query", "error", err)
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.IsActive)
		if err != nil {
			slog.Error("failed to scan user row", "error", err, "operation", "GetAll")
			return nil, err
		}
		users = append(users, user)
	}
	slog.Debug("retrieved users", "count", len(users))
	return users, nil

}

func (r *PostgresRepository) GetByID(id int) (*model.User, error) {
	query := `select id,username,email,password,name,isactive from Users where id=$1`
	row := r.db.QueryRow(query, id)
	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Name, &user.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("user not found", "user_id", id)
			return nil, errors.NewNotFoundError(id, "no user with that id")
		}
		slog.Error("error while scaning user by id", "error", err, "user_id", id)
		return nil, err
	}
	slog.Debug("user retrieved succesfully", "user_id", id)
	return &user, nil
}
func (r *PostgresRepository) GetByEmail(email string) (*model.User, error) {
	if r.ExistsByEmail(email) {
		query := `select id,username,email,name,isactive,password from Users where email=$1`
		row := r.db.QueryRow(query, email)
		var user model.User
		err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.IsActive, &user.Password)
		if err != nil {
			slog.Error("error while scanining user by email", "error", err, "user_email", email)
			return nil, err
		}
		return &user, nil
	}
	slog.Debug("user retrieved succesfully", "user_email", email)
	return nil, errors.NewNotFoundError(email, "no user with that id")

}
func (r *PostgresRepository) Create(user *model.User) error {
	query := `insert into Users (username,email,password,name) values($1,$2,$3,$4) returning id`
	err := r.db.QueryRow(query, user.Username, user.Email, user.Password, user.Name).Scan(&user.ID)
	if err != nil {
		slog.Error("failed to create user", "error", err, "user_id", user.ID, "user_email", user.Email)
		return err
	}
	slog.Info("user created successfully", "user_id", user.ID, "user_email", user.Email)
	return nil
}
func (r *PostgresRepository) Update(id int, user model.User) error {
	//can optimize by calling existsbyid here clear redundant code.
	if r.ExistsByID(id) {
		query := `update Users set username=$1 ,email=$2,password=$3,updated_at=CURRENT_TIMESTAMP  where id=$4 `
		result, err := r.db.Exec(query, user.Username, user.Email, user.Password, id)
		if err != nil {
			slog.Error("unable to execute update query", "error", err, "user_id", id)
			return fmt.Errorf("unable to exec query %w", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			slog.Error("unable to get rows affected", "error", err)
			return fmt.Errorf("unable to get rows affectted")
		}
		if rowsAffected == 0 {
			slog.Error("user not found", "user_id", id)
			return errors.NewNotFoundError(id, "no user with the given id")
		}
		return nil
	}
	slog.Info("user updated successfully", "user_id", id)
	return errors.NewNotFoundError(id, "user not found")
}
func (r *PostgresRepository) Delete(id int) error {
	//can optimize by calling existsbyid here clear redundant code.
	query := `delete from Users where id=$1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		slog.Error("unalbe to execute delete query", "error", err, "user_id", id)
		return fmt.Errorf("unable to exec query")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("unable to et rows affected", "error", err)
		return fmt.Errorf("unable to get rows affectted")
	}
	if rowsAffected == 0 {
		slog.Error("unable to get rows affected", "error", err)
		return errors.NewNotFoundError(id, "no user with the given id")
	}
	slog.Info("user deleted successfully","user_id",id)
	return nil
}

func (r *PostgresRepository) ExistsByEmail(email string) bool { // Returns bool, not error
	query := `SELECT EXISTS(SELECT 1 FROM Users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		slog.Error("failed to check email existence", "error", err)
		return false // On error, assume doesn't exist
	}
    slog.Debug("email existence check", "email", email, "exists", exists)
	return exists
}
func (r *PostgresRepository) ExistsByID(id int) bool { // Returns bool, not error
	query := `SELECT EXISTS(SELECT 1 FROM Users WHERE id = $1)`

	var exists bool
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		slog.Error("failed to check id existence", "error", err)
		return false // On error, assume doesn't exist
	}
	 slog.Debug("username existence check", "id", id, "exists", exists)
	return exists
}
func (r *PostgresRepository) ExistsByUsername(username string) bool {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	var exists bool
	err := r.db.QueryRow(query, username).Scan(&exists)
	if err!=nil{
		slog.Error("failed to check email existence", "error", err)
		return false
	}
	 slog.Debug("username existence check", "username", username, "exists", exists)
	return exists
}
