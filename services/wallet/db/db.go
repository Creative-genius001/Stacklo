package db

import (
	"context"
	"fmt"
	"os"

	"github.com/Creative-genius001/go-logger"
	"github.com/jackc/pgx/v5"
)

func TestDBConn(url string) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	logger.Info("Connection to database url successful")

	// sqlScript, err := os.ReadFile("create_wallets_table.sql")
	// if err != nil {
	// 	logger.Fatal("Unable to read SQL script file: ", err)
	// 	os.Exit(1)
	// }

	// _, err = conn.Exec(context.Background(), string(sqlScript))
	// if err != nil {
	// 	logger.Fatal("Failed to execute SQL script: ", err)
	// 	os.Exit(1)
	// }

	// logger.Info("Wallets table created or already exists successfully!")

}
