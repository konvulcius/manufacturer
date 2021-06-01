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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
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

type instrumentingMiddleware struct {
	CliDB      *sqlx.DB
	HTTPClient *http.Client
	Latencies  *prometheus.HistogramVec
}

func (i *instrumentingMiddleware) getManufacturer(ctx context.Context, n int) (result []*manufacturer, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	defer i.recordMetrics("postgres_get", time.Now())

	err = i.CliDB.SelectContext(ctx, &result, `SELECT m.id, m.details FROM manufacture.manufacturer m
	WHERE (details ->> 'needUpdate')::boolean = true LIMIT $1`, n)

	return
}

func (i *instrumentingMiddleware) doRequest(ctx context.Context, manufacturers []*manufacturer) error {
	data, err := json.Marshal(manufacturers)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/api/v1/manufacturer", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	req = req.WithContext(ctx)

	defer i.recordMetrics("api", time.Now())

	resp, err := i.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New("status code is error")
	}

	return nil
}

func (i *instrumentingMiddleware) updateManufacturer(ctx context.Context, IDs []string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	query, args, err := sqlx.In(`UPDATE manufacture.manufacturer SET details = details|| '{"needUpdate": false}'::jsonb
WHERE id IN (SELECT unnest(array[?]))`, IDs)
	if err != nil {
		return err
	}
	query = i.CliDB.Rebind(query)

	defer i.recordMetrics("postgres_update", time.Now())

	_, err = i.CliDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (i *instrumentingMiddleware) recordMetrics(metric string, startTime time.Time) {
	labels := map[string]string{
		"name": metric,
	}
	i.Latencies.With(labels).Observe(time.Since(startTime).Seconds())
}

func newInstrumentingMiddleware(cliDB *sqlx.DB, httpClient *http.Client) *instrumentingMiddleware {
	return &instrumentingMiddleware{
		CliDB:      cliDB,
		HTTPClient: httpClient,
		Latencies: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "script_api_latency",
				Help: "Requests Latency",
				Buckets: []float64{
					0.01, 0.05, 0.100,
					0.5, 0.95, 1,
				},
			},
			[]string{"name"},
		),
	}
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

	instruments := newInstrumentingMiddleware(db, client)
	prometheus.MustRegister(instruments.Latencies)

	defer func() {
		if err := push.New("http://127.0.0.1:9091", "script_latencies").Collector(instruments.Latencies).Push(); err != nil {
			_ = level.Error(logger).Log("msg", "failed to push metrics", "err", err)
		}
	}()

	for {
		// take data from db
		manufacturers, err := instruments.getManufacturer(ctx, *n)
		if err != nil {
			_ = level.Error(logger).Log("msg", "failed connect to get data from postgres", "err", err)
			os.Exit(1)
		}
		if manufacturers == nil {
			break
		}

		// do request to api
		err = instruments.doRequest(ctx, manufacturers)
		if err != nil {
			_ = level.Error(logger).Log("msg", "failed to do request", "err", err)
			os.Exit(1)
		}

		// do update
		unnestArray := make([]string, 0, len(manufacturers))
		for _, man := range manufacturers {
			unnestArray = append(unnestArray, man.ID)
		}

		err = instruments.updateManufacturer(ctx, unnestArray)
		if err != nil {
			_ = level.Error(logger).Log("msg", "failed to update manufacturers", "err", err)
			os.Exit(1)
		}
	}

	_ = logger.Log("msg", "Well Done")
}
