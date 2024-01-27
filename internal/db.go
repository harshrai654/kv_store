package internal

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB
const DEFAULT_MAX_OPEN_CONNECTIONS = 10
const DEFAULT_MAX_IDLE_CONNECTIONS = 5


func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	maxOpenConnections, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONN"))
	if err != nil {
		maxOpenConnections = DEFAULT_MAX_OPEN_CONNECTIONS
	}

	maxIdleConnections, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONN"))
	if err != nil {
		maxIdleConnections = DEFAULT_MAX_IDLE_CONNECTIONS
	}

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword,dbHost,dbPort,dbName)

	db, err = sql.Open("mysql", dataSourceName)

	if err != nil {
		log.Fatal("DB connection failed", err)
	}

	//Configure connection pool
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetMaxOpenConns(maxOpenConnections)

	err = db.Ping()

	if err != nil {
		log.Fatal("Unable to communicate to DB ", err)
	}

	log.Println("DB Connection successful!")
}


func GetValueFromKey(key string) []byte {
	row := db.QueryRow(fmt.Sprintf("SELECT v FROM kv WHERE k = '%s' AND (expired_at IS NULL OR expired_at > NOW())", key))
	var value []byte

	err := row.Scan(&value)

	if err != nil {
		return nil;
	}
	return value
}

func UpdateKey(key string, value string) (bool, int64) {
	var query string
	var result sql.Result
	var err error

	query = `UPDATE kv SET v = ? WHERE k = ? AND expired_at > NOW()`

	result, err = db.Exec(query, value, key)

	if err != nil {
		log.Println(err)
		return false,0
	}

	var rowCount int64
	rowCount, err = result.RowsAffected()

	if err != nil {
		log.Println(err)
		return false,0
	}

	
	log.Printf("Update complete. Rows affected: %d", rowCount)
	
	if rowCount == 0 {
		return true, 0
	}

	return true, rowCount
}

func InsertKey(key string, value string, expired_at int64) (bool, int64) {
	var query string
	var result sql.Result
	var err error

	if expired_at != 0 {
		query = `INSERT INTO kv (k, v, expired_at) VALUES (?,?,FROM_UNIXTIME(?)) 
		ON DUPLICATE KEY UPDATE v = ?, expired_at = FROM_UNIXTIME(?)`
		result, err = db.Exec(query, key, value, expired_at, value, expired_at)
	} else {
		query = `INSERT INTO kv (k, v) VALUES (?,?) ON DUPLICATE KEY 
		UPDATE v = ?, 
		expired_at = CASE WHEN expired_at > NOW() THEN expired_at ELSE NULL END;`
		result, err = db.Exec(query, key, value, value)
	}


	if err != nil {
		log.Println(err)
		return false,0
	}

	var rowCount int64
	rowCount, err = result.RowsAffected()

	if err != nil {
		log.Println(err)
		return false,0
	}

	
	log.Printf("Insert complete. Rows affected: %d", rowCount)
	
	if rowCount == 0 {
		return true, 0
	}

	return true, rowCount
}

func DeleteKey(key string) (bool, int64){
	result, err := db.Exec("UPDATE kv SET expired_at = NOW() - 1 WHERE k = ? AND expired_at > NOW()", key)

	if err != nil {
		log.Println(err)
		return false, 0;
	}

	var rowCount int64
	rowCount, err = result.RowsAffected()

	if err != nil {
		log.Println(err)
		return false, 0;
	}

	log.Printf("Key added for delete job, Rows Affected: %d", rowCount)

	if rowCount == 0 {
		return true, 0
	}

	return true, rowCount
}

func UpdateTTL(key string, expired_at int64) (bool, int64) {
	result, err := db.Exec("UPDATE kv SET expired_at = FROM_UNIXTIME(?) WHERE k = ? AND (expired_at IS NULL OR expired_at > NOW())", expired_at, key)

	if err != nil {
		log.Println(err)
		return false, 0;
	}

	var rowCount int64
	rowCount, err = result.RowsAffected()

	if err != nil {
		log.Println(err)
		return false, 0;
	}

	log.Printf("TTL updated, Rows Affected: %d", rowCount)

	if rowCount == 0 {
		return true, 0
	}

	return true, rowCount
}