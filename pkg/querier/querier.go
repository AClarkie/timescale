package querier

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/AClarkie/timescale/pkg/csv"
	"go.uber.org/zap"
)

var (
	minute    time.Time
	host      string
	max_usage float64
	min_usage float64
)

type querier struct {
	Hostnames []csv.Hostname
	Logger    *zap.SugaredLogger
	Statement string
	DB        *sql.DB
}

func NewQuerier(hostname []csv.Hostname, logger *zap.SugaredLogger, db *sql.DB) (*querier, error) {
	return &querier{
		Hostnames: hostname,
		Logger:    logger,
		DB:        db,
	}, nil
}

func (q *querier) Start() error {
	fmt.Println("Starting querier")

	for _, hostname := range q.Hostnames {
		for _, timeRange := range hostname.Times {

			start := time.Now()
			rows, err := q.DB.Query(`SELECT
        time_bucket('1 minute', ts) AS minute,
        host,
        max(usage) AS max_usage,
        min(usage) AS min_usage
      FROM cpu_usage
      WHERE ts >= $2 AND ts < $3
      AND host = $1
      GROUP BY minute, host
      ORDER BY minute DESC;`, hostname.Name, timeRange.StartTime, timeRange.EndTime)

			if err != nil {
				q.Logger.Fatal(err)
			}
			elapsed := time.Since(start)
			fmt.Println(elapsed)
			defer rows.Close()

			// Loop through rows, using Scan to assign column data to struct fields.
			for rows.Next() {
				err := rows.Scan(&minute, &host, &max_usage, &min_usage)
				if err != nil {
					log.Fatal(err)
				}
				// q.Logger.Info(minute, host, max_usage, min_usage)
			}

		}
	}

	return nil
}
