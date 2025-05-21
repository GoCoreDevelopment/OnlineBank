package repository

import (
	"api/internal/config"
	"api/internal/db"
	"api/internal/models/user"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)


func setupDB(t *testing.T) *sql.Tx {
	cfg := config.Config{
		DataBaseURL: "postgres://postgres:1234@localhost:5432/postgres?sslmode=disable",
	}
	conn, err := db.InitDB(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	tx, err := conn.Begin()
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	return tx
}

func TestRepository_TransferMoney(t *testing.T) {
	dataUsers := []user.UserRegistred{
		{
			FirstName: "Ivan",
			LastName:  "Bubnov",
			Phone:     "+79999999999",
			Email:     "bubnov.ivanbubnoff@yandex.ru",
			Password:  "1234",
		},
		{
			FirstName: "Sergei",
			LastName:  "Bubnov",
			Phone:     "+71111111111",
			Email:     "serg.bubnov@yandex.ru",
			Password:  "1111",
		},
	}

	conn := setupDB(t)
	defer conn.Rollback()

	repo := NewRepository(conn)

	var senderID, receiverID int
	for i, v := range dataUsers {
		id, err := repo.CreateUser(v)
		assert.NoError(t, err)
		if i == 0 {
			senderID = id
		} else {
			receiverID = id
		}

		_, err = conn.Exec("UPDATE users SET balance = 1000 WHERE id = $1", id)
		assert.NoError(t, err)
	}

	t.Run("Successful Transfer", func(t *testing.T) {
		err := repo.TransferMoney(senderID, receiverID, 250)
		assert.NoError(t, err)

		balanceSender, err := repo.GetUserBalance(senderID)
		assert.NoError(t, err)
		assert.Equal(t, 750, balanceSender)

		balanceReceiver, err := repo.GetUserBalance(receiverID)
		assert.NoError(t, err)
		assert.Equal(t, 1250, balanceReceiver)
	})

	t.Run("Insufficient Balance", func(t *testing.T) {
		err := repo.TransferMoney(senderID, receiverID, 1000) // Баланс отправителя = 750
		assert.Error(t, err)
		assert.Equal(t, "insufficient balance", err.Error())
	})

	t.Run("Receiver Not Found", func(t *testing.T) {
		err := repo.TransferMoney(senderID, 999, 100) // ID 999 не существует
		assert.Error(t, err)
		assert.Equal(t, "receiver not found", err.Error())
	})

	t.Run("Sender Not Found", func(t *testing.T) {
		err := repo.TransferMoney(999, receiverID, 100) // ID 999 не существует
		assert.Error(t, err)
		assert.Equal(t, "sender not found", err.Error())
	})
}