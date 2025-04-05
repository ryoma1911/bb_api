package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// 仮の環境変数ファイルを作成
func createInvalidEnvFile(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "invalid_*.env")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(content)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// TestGetDSNFromEnv： 環境変数「.env」ファイルから各パラメータを取得するケースをテスト
func TestGetDSNFromEnv(t *testing.T) {
	connect := DBService{}

	t.Run("Valid envfile", func(t *testing.T) {
		//DSNが正しく取得されていること
		dsn, err := connect.GetDSNFromEnv("/code/.env")
		assert.NoError(t, err)

		// 期待値のDSNを定義
		expected := "bbapi:bbapi@tcp(bb_db:3306)/bbapi-db"

		//返却結果が期待通りであることを確認
		assert.Equal(t, dsn, expected)

	})

	t.Run("Invalid envfile", func(t *testing.T) {
		//無功な環境変数ファイルでエラーが発生すること
		_, err := connect.GetDSNFromEnv("/.ttt")
		assert.Error(t, err)
	})

	t.Run("Invalid envparam", func(t *testing.T) {
		//テスト用に仮ファイルを作成
		content := `INVALID_LINE
		=NO_KEY
		FOO`

		filename, err := createInvalidEnvFile(content)
		assert.NoError(t, err)
		defer os.Remove(filename)

		//無功な環境変数の取得でエラーが発生すること
		_, err = connect.GetDSNFromEnv(filename)
		assert.Error(t, err)
	})

}

// TestConnectOnly： DB接続するケースをテスト
func TestConnectOnly(t *testing.T) {
	connect := DBService{}

	t.Run("Valid DSN", func(t *testing.T) {
		//設定したDSNでDB接続に成功していること
		dsn := "bbapi:bbapi@tcp(bb_db:3306)/bbapi-db"

		_, err := connect.ConnectOnly(dsn)
		assert.NoError(t, err)
	})

	t.Run("Invalid DSN", func(t *testing.T) {
		//無功なDSN接続でエラーが発生する事
		dsn := "ttt"
		_, err := connect.ConnectOnly(dsn)
		assert.Error(t, err)
	})

}

// Testcheckconnect： DB接続後のpingを送信するケースをテスト
func TestCheckconnect(t *testing.T) {
	// SQL モックの作成
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	t.Run("Valid ping", func(t *testing.T) {
		// ping が通るモックを作成
		mock.ExpectPing()

		conn, err := checkconnect(db)
		assert.NoError(t, err)

		if conn == nil {
			t.Fatal("checkconnect returned nil db")
		}

		// 期待された操作がすべて満たされたか確認
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("Invalid ping", func(t *testing.T) {
		// ping エラーのモックを作成
		mock.ExpectPing().WillReturnError(fmt.Errorf("mock ping error"))
		mock.ExpectClose()

		conn, err := checkconnect(db)
		assert.Error(t, err)

		if conn != nil {
			t.Fatal("Expected nil db on error but got non-nil db")
		}

		// 期待された操作がすべて満たされたか確認
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}
