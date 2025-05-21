package repository

import (
	"database/sql"
	"errors"

	"api/internal/models/user"
)

type DBTX interface {
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}	

type Repository interface {
	CreateUser(user user.UserRegistred) (int, error)
	GetUserPassword(email string) (string, error)
	GetUserBalance(id int) (int, error)
	GetUserID(emain string) (int, error)
	TransferMoney(senderID, receiverID, amount int) error
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(user user.UserRegistred) (int, error) {
	query := "INSERT INTO users (first_name, last_name, phone, email, password) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	var id int
	err := r.db.QueryRow(query, user.FirstName, user.LastName, user.Phone, user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) GetUserPassword(email string) (string, error) {
	query := "SELECT id, password FROM users WHERE email=$1"

	var id int
	var password string
	err := r.db.QueryRow(query, email).Scan(&id, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("the user won't find it")
		}
		return "", err
	}

	return password, nil
}

func (r *repository) GetUserBalance(id int) (int, error) {
	query := "SELECT balance FROM users WHERE id=$1"

	var balance int
	err := r.db.QueryRow(query, id).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (r *repository) GetUserID(email string) (int, error) {
	query := "SELECT id FROM users WHERE email=$1"

	var id int
	err := r.db.QueryRow(query, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("the user won't find it")
		}
		return 0, err
	}

	return id, nil
}

func (r *repository) TransferMoney(senderID, receiverID, amount int) (err error) {
	tx, ok := r.db.(*sql.Tx)
	if !ok {
		tx, err = r.db.(*sql.DB).Begin()
		if err != nil {
			return err
		}

		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}

	var senderBalance int
	query := "SELECT balance FROM users WHERE id=$1 FOR UPDATE"
	err = tx.QueryRow(query, senderID).Scan(&senderBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("sender not found")
		}
		return err
	}

	if senderBalance < amount {
		return errors.New("insufficient balance")
	}

	var receiverExist bool
	query = "SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)"
	err = tx.QueryRow(query, receiverID).Scan(&receiverExist)
	if err != nil {
		return err
	}

	if !receiverExist {
		return errors.New("receiver not found")
	}

	query = "UPDATE users SET balance = balance - $1 WHERE id=$2"
	_, err = tx.Exec(query, amount, senderID)
	if err != nil {
		return err
	}

	query = "UPDATE users SET balance = balance + $1 WHERE id=$2"
	_, err = tx.Exec(query, amount, receiverID)
	if err != nil {
		return err
	}

	return nil
}
