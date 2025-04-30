package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestConnectOnly(t *testing.T) {
	connect := DBService{}

	t.Run("Valid DSN from env", func(t *testing.T) {
		// テスト用の環境変数を設定
		_ = godotenv.Load("/code/.env")

		// DSNの各環境変数が存在することを確認
		assert.NotEmpty(t, os.Getenv("MYSQL_USER"))
		assert.NotEmpty(t, os.Getenv("MYSQL_PASSWORD"))
		assert.NotEmpty(t, os.Getenv("MYSQL_DATABASE"))

		db, err := connect.ConnectOnly()
		assert.NoError(t, err)

		if db != nil {
			_ = db.Close()
		}
	})

	t.Run("Invalid env settings", func(t *testing.T) {
		// 不正な環境変数に書き換える
		os.Setenv("MYSQL_USER", "")
		os.Setenv("MYSQL_PASSWORD", "")
		os.Setenv("MYSQL_DATABASE", "")

		_, err := connect.ConnectOnly()
		assert.Error(t, err)
	})
}

func TestCheckconnect(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	t.Run("Valid ping", func(t *testing.T) {
		mock.ExpectPing()

		conn, err := checkconnect(db)
		assert.NoError(t, err)
		assert.NotNil(t, conn)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Invalid ping", func(t *testing.T) {
		mock.ExpectPing().WillReturnError(fmt.Errorf("mock ping error"))
		mock.ExpectClose()

		conn, err := checkconnect(db)
		assert.Error(t, err)
		assert.Nil(t, conn)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestConnect(t *testing.T) {
	getDSN()
}
