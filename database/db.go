package database

import (
	"fmt"
	"strconv"
	"strings"

	// source/file import is required for migration files to read
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	// load pq as database driver
	_ "github.com/lib/pq"
)

var (
	Data *sqlx.DB
)

type SSLMode string

const (
	SSLModeEnable  SSLMode = "enable"
	SSLModeDisable SSLMode = "disable"
)

// ConnectAndMigrate function connects with a given database and returns error if there is any error
func ConnectAndMigrate(host string, port int, databaseName, user, password string, sslMode SSLMode) *sqlx.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, databaseName, sslMode)

	// DB, err := sqlx.Open("postgres", connStr)
	DB, err := sqlx.Open("postgres", connStr)

	if err != nil {
		return nil
	}

	err = DB.Ping()
	if err != nil {
		return nil
	}
	Data = DB
	err = migrateUp(DB)
	if err != nil {
		return nil
	}
	return DB
}

func ShutdownDatabase() error {
	return Data.Close()
}

// migrateUp function migrate the database and handles the migration logic
func migrateUp(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"mydb", driver)

	if err != nil {
		logrus.Errorf("error in migration: %s", err)
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// Tx provides the transaction wrapper
func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := Data.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %+v", err)
	}
	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				logrus.Errorf("failed to rollback tx: %s", rollBackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			logrus.Errorf("failed to commit tx: %s", commitErr)
		}
	}()
	err = fn(tx)
	return err
}

// SetupBindVars prepares the SQL statement for batch insert
func SetupBindVars(stmt, bindVars string, length int) string {
	bindVars += ","
	stmt = fmt.Sprintf(stmt, strings.Repeat(bindVars, length))
	return replaceSQL(strings.TrimSuffix(stmt, ","), "?")
}

// replaceSQL replaces the instance occurrence of any string pattern with an increasing $n based sequence
func replaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}
