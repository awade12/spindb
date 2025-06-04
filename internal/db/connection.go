package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

type ConnectionTester struct{}

func NewConnectionTester() *ConnectionTester {
	return &ConnectionTester{}
}

func (ct *ConnectionTester) TestPostgres(host string, port int, user, password, dbname string) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return ct.testConnection("postgres", dsn)
}

func (ct *ConnectionTester) TestMySQL(host string, port int, user, password, dbname string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname)

	return ct.testConnection("mysql", dsn)
}

func (ct *ConnectionTester) TestSQLite(filePath string) error {
	return ct.testConnection("sqlite", filePath)
}

func (ct *ConnectionTester) testConnection(driver, dsn string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(5 * time.Second)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func (ct *ConnectionTester) WaitForDatabase(driver, dsn string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if err := ct.testConnection(driver, dsn); err == nil {
			return nil
		}
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("database did not become available within %v", timeout)
}
