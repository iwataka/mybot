package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func initDB(c *DBConfig) error {
	if c == nil {
		return nil
	}
	if db == nil {
		driver := c.Driver
		dataSource := c.DataSource
		if len(driver) == 0 || len(dataSource) == 0 {
			return nil
		}
		var err error
		db, err = sql.Open(driver, dataSource)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertVisionDBColumn(c *DBConfig, imageURL, result string) error {
	err := initDB(c)
	if err != nil {
		return err
	}
	if db == nil {
		return nil
	}
	table := c.VisionTable
	if len(table) == 0 {
		return nil
	}
	cmd := fmt.Sprintf("INSERT INTO %s VALUES ('%s', '%s')", table, imageURL, result)
	_, err = db.Exec(cmd)
	if err != nil {
		return err
	}
	return nil
}

type VisionDBColumn struct {
	ImageURL string
	Result   string
}

func selectVisionDBColumn(c *DBConfig, num int) ([]*VisionDBColumn, error) {
	err := initDB(c)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return make([]*VisionDBColumn, 0, 0), nil
	}
	table := c.VisionTable
	if len(table) == 0 {
		return make([]*VisionDBColumn, 0, 0), nil
	}
	cmd := fmt.Sprintf("SELECT LAST(%d) * FROM %s", num, table)
	rows, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}
	result := make([]*VisionDBColumn, num, num)
	for rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		item := &VisionDBColumn{cols[0], cols[1]}
		result = append(result, item)
	}
	return result, nil
}
