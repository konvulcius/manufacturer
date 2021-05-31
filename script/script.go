package main

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type manufacturer struct {
	ID      string  `json:"id" db:"id"`
	Details Details `json:"details" db:"details"`
}

type Details map[string]interface{}

func (a *Details) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Details) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func connect(ctx context.Context) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	db, err := sqlx.ConnectContext(ctx, "pgx", "postgresql://postgres:postgres@127.0.0.1:54320/mydb")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	return db, nil
}

func main() {
	n := flag.Int("n", 3, "ограничение на количество считывания данных")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestamp)

	db, err := connect(ctx)
	if err != nil {
		_ = level.Error(logger).Log("msg", "failed connect to postgresql", "err", err)
	}

	client := &http.Client{}

	for {
		// take data from db
		var manufacturers []*manufacturer
		err = db.SelectContext(ctx, &manufacturers, `SELECT m.id, m.details FROM manufacture.manufacturer m
	WHERE (details ->> 'needUpdate')::boolean = true LIMIT $1`, n)
		if err != nil {
			_ = level.Error(logger).Log("msg", "failed connect to get data from postgres", "err", err)
			os.Exit(1)
		}
		if manufacturers == nil {
			break
		}

		// do request to api
		data, err := json.Marshal(manufacturers)
		if err != nil {
			_ = level.Error(logger).Log("msg", "failed to marshal data", "err", err)
			os.Exit(1)
		}
		req, err := http.NewRequest("POST", "http://127.0.0.1:8000/api/v1/manufacturer", bytes.NewBuffer(data))
		if err != nil {
			_ = level.Error(logger).Log("msg", "failed to make request", "err", err)
			os.Exit(1)
		}
		req.Header.Set("Content-Type", "application/json")

		req = req.WithContext(ctx)

		resp, err := client.Do(req)
		if err != nil {
			_ = level.Error(logger).Log("msg", "failed to do request", "err", err)
			os.Exit(1)
		}

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			_ = level.Error(logger).Log("msg", "status code is error", "err", err)
			os.Exit(1)
		}
		resp.Body.Close()

		// do update
		unnestArray := make([]string, 0, len(manufacturers))
		for _, man := range manufacturers {
			unnestArray = append(unnestArray, man.ID)
		}

		query, args, err := sqlx.In(`UPDATE manufacture.manufacturer SET details = details|| '{"needUpdate": false}'::jsonb
WHERE id IN (SELECT unnest(array[?]))`, unnestArray)
		if err != nil {
			_ = level.Error(logger).Log("msg", "error", "err", err)
			os.Exit(1)
		}
		query = db.Rebind(query)

		_, err = db.ExecContext(ctx, query, args...)
		if err != nil {
			_ = level.Error(logger).Log("msg", "failed to update data in postgres", "err", err)
			os.Exit(1)
		}
	}

	_ = logger.Log("msg", "Well Done")
}
