package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// DBHandler インターフェースで関数を抽象化
type DBHandler interface {
	GetDSNFromEnv(file string) (string, error)
	ConnectOnly(dsn string) (*sql.DB, error)
}

// DBService デフォルトの実装
type DBService struct{}

// DSNを.envから読み取って生成
func (d *DBService) GetDSNFromEnv(file string) (string, error) {
	// .envを読み込み
	err := godotenv.Load(file)
	if err != nil {
		return "", fmt.Errorf("failed to load env file: %w", err)
	}

	// 環境変数からDB接続情報を取得
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")

	// DSNを生成
	dsn := fmt.Sprintf("%s:%s@tcp(bb_db:3306)/%s", user, password, dbName)

	return dsn, nil
}

func checkconnect(db *sql.DB) (*sql.DB, error) {
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to check to connect database: %w", err)
	}
	log.Println("success connect to database")
	return db, nil
}

func (d *DBService) ConnectOnly(dsn string) (*sql.DB, error) {

	//データベースのハンドルを取得
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	return checkconnect(db)
}
