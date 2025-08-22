package repository

import (
	"database/sql"
	"fmt"
	"user-management/internal/model"
)

type UserRepo interface{
	GetAll()([]model.User,error)
	GetByID()(*model.User,error)
	Create(user model.User)error
	Update(id int,user model.User)error
	Delete(id int)error
}

type PostgresRepository struct{
	db *sql.DB
}

func NewPostgresRepository(connectionString string) (UserRepo,error){
	db,err:=sql.Open("postgres",connectionString)
	if err!=nil{
		return nil,fmt.Errorf("unable to open db conection:%w",err)
	}
	return &PostgresRepository{db: db},nil
}

func(r *PostgresRepository)GetAll()([]model.User,error){
	query:=`select id,username,email,name,isactive from Users`
	rows,err:=r.db.Query(query)
	if err!=nil{
		return nil, fmt.Errorf("error while querying:%w",err)

	}
	defer rows.Close()
	var users []model.User
	for rows.Next(){
		var user model.User
		err:=rows.Scan(&user.ID,&user.Username,&user.Email,&user.Name,&user.IsActive)
		if err!=nil{
			return nil,fmt.Errorf("error while scaning row data:%w ",err)
		}
		users=append(users, user)
	}
	return users,nil

}

func(r *PostgresRepository)GetByID()(*model.User,error){
	query:=`select * from Users`
	row:=r.db.QueryRow(query)
	var user model.User
	err:=row.Scan(&user.ID,&user.Username,&user.Email,&user.Password,&user.Name,&user.IsActive)
	if err!=nil{
			return nil,fmt.Errorf("error while scaning row data:%w ",err)
		}
	return &user,nil
	}
func(r *PostgresRepository)	Create(user model.User)error{
	query:=`insert into Users (username,email,password) values($1,$2,$3) returning id`
	err:=r.db.QueryRow(query,user.Username,user.Email,user.Password).Scan(&user.ID)
	if err!=nil{
		return fmt.Errorf("error while scaning row data:%w ",err)
		}
	return nil
}
func(r *PostgresRepository)	Update(id int,user model.User)error{
	query:=`update Users set username=$1 email=$2 password=$3 updated_at=CURRENT_TIMESTAMP  where id=$4 `
	result,err:=r.db.Exec(query,user.Username,user.Email,user.Password,user.ID)
	if err!=nil{
		return fmt.Errorf("unable to exec query")
	}
	rowsAffected,err:=result.RowsAffected()
	if err!=nil{
		return fmt.Errorf("unable to get rows affectted")
	}
	if rowsAffected==0{
		return fmt.Errorf("no such user with the provided id")
	}
	return nil
}
func(r *PostgresRepository)	Delete(id int)error{
	query:=`delete from Users where id=$1`
	result,err:=r.db.Exec(query,id)
	if err!=nil{
		return fmt.Errorf("unable to exec query")
	}
	rowsAffected,err:=result.RowsAffected()
	if err!=nil{
		return fmt.Errorf("unable to get rows affectted")
	}
	if rowsAffected==0{
		return fmt.Errorf("no such user with the provided id")
	}
	return nil
}