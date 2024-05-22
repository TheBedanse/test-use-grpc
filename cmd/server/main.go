package main

import (
	"context"
	"fmt"
	connectDB "grpc/cmd/connectdb"
	"grpc/pkg/api"
	chart "grpc/pkg/charts"
	"grpc/pkg/handlers"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"google.golang.org/grpc"
)

func main() {
	conn, err := connectDB.СonnectToDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error DB: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	s := grpc.NewServer()
	api.RegisterCountServer(s, handlers.Handler)

	go func() {
		l, err := net.Listen("tcp", ":8000")
		if err != nil {
			log.Fatal(err)
		}
		if err := s.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()
	var count int = 0
	var change int = 0
	var countid int = 1
	var createdat = time.Now()

	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		rows, err := conn.Query(context.Background(), "select count_id, created_at, current_value, change from counter_history where is_deleted = false")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
			os.Exit(1)
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&countid, &createdat, &count, &change)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
				return
			}
			str := fmt.Sprintf("<h1>ID: %d Время: %s Знач: %d Изменилось на: %d</h1>", countid, createdat.Format(time.RFC3339), count, change)
			fmt.Fprint(w, str)
		}
	})

	http.HandleFunc("/chart", func(w http.ResponseWriter, r *http.Request) {
		line := charts.NewLine()

		line.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
			charts.WithTitleOpts(opts.Title{
				Title: "Last 10 score in counter",
			}))

		line.SetXAxis([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}).
			AddSeries("Scores1", chart.GenerateLineItems(1)).
			AddSeries("Scores2", chart.GenerateLineItems(2)).
			AddSeries("Scores3", chart.GenerateLineItems(3)).
			AddSeries("Scores4", chart.GenerateLineItems(4)).
			AddSeries("Scores5", chart.GenerateLineItems(5)).
			SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: false}))
		line.Render(w)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
