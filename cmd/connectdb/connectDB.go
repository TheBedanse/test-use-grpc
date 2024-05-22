package connectDB

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func СonnectToDB() (*pgx.Conn, error) {
	urlExample := "postgres://admin:root@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error DB: %v\n", err)
		os.Exit(1)
	}
	return conn, nil
}
