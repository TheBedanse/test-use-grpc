package handlers

import (
	"context"
	"fmt"
	connectDB "grpc/cmd/connectdb"
	"grpc/pkg/api"
	"os"
	"time"
)

type Service struct {
	api.UnimplementedCountServer
}

var Handler Service

func (s Service) Change(ctx context.Context, req *api.ChangeRequest) (*api.CountResponse, error) {
	conn, err := connectDB.СonnectToDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error DB: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	var count int = 0
	var change int = 0
	var countid int = int(req.GetCountid())
	var createdat = time.Now()

	err = conn.QueryRow(context.Background(), "select current_value from counters_more where count_id = $1", countid).Scan(&count)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}
	change = int(req.GetX())
	count = count + change
	_, err = conn.Exec(context.Background(), "insert into counter_history (count_id, created_at, current_value, change) values ($1, $2, $3, $4)",
		countid, createdat, count, change)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exec fail: %v\n", err)
		os.Exit(1)
	}
	_, err = conn.Exec(context.Background(), "update counters_more set current_value = $1 where count_id = $2 and is_deleted = false",
		count, countid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exec fail: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(count, change)
	return &api.CountResponse{Result: int32(count)}, nil
}

func (s Service) AddCounter(ctx context.Context, req *api.NameCountRequest) (*api.ReadyCountResponse, error) {
	conn, err := connectDB.СonnectToDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error DB: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	var count int = 0
	var change int = 0
	var countid int
	var createdat = time.Now()
	var countname string = req.GetName()
	if req.GetName() == "" {
		countname = "New counter"
	}
	_, err = conn.Exec(context.Background(), "insert into counters_more (name_count, current_value) values ($1, $2)",
		countname, count)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exec fail: %v\n", err)
		os.Exit(1)
	}
	err = conn.QueryRow(context.Background(), "select count_id from counters_more where count_id = (select max(count_id) from counters_more)").Scan(&countid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}
	_, err = conn.Exec(context.Background(), "insert into counter_history (count_id, created_at, current_value, change) values ($1, $2, $3, $4)",
		countid, createdat, count, change)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exec fail: %v\n", err)
		os.Exit(1)
	}

	return &api.ReadyCountResponse{Name: countname}, nil
}

func (s Service) Score(ctx context.Context, req *api.CountRequest) (*api.CountResponse, error) {
	conn, err := connectDB.СonnectToDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error DB: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())
	var countid int = int(req.GetCountid())
	var count int
	err = conn.QueryRow(context.Background(), "select current_value from counters_more where count_id = $1 and is_deleted = false", countid).Scan(&count)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}

	return &api.CountResponse{Result: int32(count)}, nil
}

func (s Service) Delete(ctx context.Context, req *api.CountRequest) (*api.ReadyDeleteResponse, error) {
	conn, err := connectDB.СonnectToDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error DB: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	var countid int = int(req.GetCountid())
	deleted := "true"

	_, err = conn.Exec(context.Background(), "update counter_history set is_deleted = true where count_id = $1",
		countid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exec fail: %v\n", err)
		os.Exit(1)
	}
	_, err = conn.Exec(context.Background(), "update counters_more set is_deleted = true where count_id = $1",
		countid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exec fail: %v\n", err)
		os.Exit(1)
	}

	return &api.ReadyDeleteResponse{CounterDelete: deleted}, nil
}
