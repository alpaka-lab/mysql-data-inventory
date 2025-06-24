package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Update with your MySQL connection string
	dsn := "user:password@tcp(localhost:3306)/information_schema"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Get all non-system databases
	dbs, err := db.Query(`
		SELECT SCHEMA_NAME
		FROM SCHEMATA
		WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
	`)
	if err != nil {
		log.Fatalf("Failed to query databases: %v", err)
	}
	defer dbs.Close()

	var databases []string
	for dbs.Next() {
		var dbname string
		if err := dbs.Scan(&dbname); err != nil {
			log.Fatalf("Failed to scan database name: %v", err)
		}
		databases = append(databases, dbname)
	}

	file, err := os.Create("mysql_inventory.csv")
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"Schema", "Table", "Column", "Type"})

	for _, schema := range databases {
		query := `
			SELECT TABLE_NAME, COLUMN_NAME, DATA_TYPE
			FROM COLUMNS
			WHERE TABLE_SCHEMA = ?
			ORDER BY TABLE_NAME, ORDINAL_POSITION
		`
		rows, err := db.Query(query, schema)
		if err != nil {
			log.Printf("Failed to query schema %s: %v", schema, err)
			continue
		}
		for rows.Next() {
			var table, column, dtype string
			if err := rows.Scan(&table, &column, &dtype); err != nil {
				log.Fatalf("Failed to scan row: %v", err)
			}
			writer.Write([]string{schema, table, column, dtype})
		}
		rows.Close()
	}

	fmt.Println("Inventory saved to mysql_inventory.csv")
}
