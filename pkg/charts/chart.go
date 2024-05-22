package chart

import (
	"context"
	"fmt"
	connectDB "grpc/cmd/connectdb"
	"os"

	"github.com/go-echarts/go-echarts/v2/opts"
)

func GenerateLineItems(countid int) []opts.LineData {
	conn, err := connectDB.СonnectToDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error DB: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	var count int
	rows, err := conn.Query(context.Background(), "select * from (select current_value from counter_history where count_id = $1 and is_deleted = false order by id asc limit 10) subquery;", countid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	items := make([]opts.LineData, 0)
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
			return nil
		}
		items = append(items, opts.LineData{Value: count})
	}
	return items
}
