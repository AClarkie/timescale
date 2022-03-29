package querier

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/AClarkie/timescale/pkg/csv"
	"github.com/AClarkie/timescale/pkg/results"
	"go.uber.org/zap"
)

func Execute(hostnames []csv.Hostname, logger *zap.SugaredLogger, db *sql.DB) *results.Result {
	var queryTimes []time.Duration
	for _, hostname := range hostnames {
		logger.Debugf("Starting queries for hostname %s", hostname.Name)
		for _, timeRange := range hostname.Times {

			start := time.Now()
			// Use Exec as we don't care about the result
			_, err := db.Exec(`SELECT
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
				logger.Debugf("Query failed for hostname %s", hostname.Name)
				return &results.Result{QueryTimes: nil, Error: err}
			}
			queryTimes = append(queryTimes, time.Since(start))
		}
	}

	logger.Debug("Finished executing queries")
	return &results.Result{QueryTimes: queryTimes, Error: nil}
}

// var (
// 	minute    time.Time
// 	host      string
// 	max_usage float64
// 	min_usage float64
// )

// func (q *querier) Start() *results.Result {
// 	var queryTimes []time.Duration
// 	for _, hostname := range q.Hostnames {
// 		for _, timeRange := range hostname.Times {

// 			start := time.Now()
// 			rows, err := q.DB.Query(`SELECT
// 					time_bucket('1 minute', ts) AS minute,
// 					host,
// 					max(usage) AS max_usage,
// 					min(usage) AS min_usage
// 				FROM cpu_usage
// 				WHERE ts >= $2 AND ts < $3
// 				AND host = $1
// 				GROUP BY minute, host
// 				ORDER BY minute DESC;`, hostname.Name, timeRange.StartTime, timeRange.EndTime)

// 			if err != nil {
// 				return &results.Result{QueryTimes: nil, Error: err}
// 			}
// 			queryTimes = append(queryTimes, time.Since(start))
// 			defer rows.Close()

// 			// Loop through rows, using Scan to assign column data to struct fields.
// 			for rows.Next() {
// 				err := rows.Scan(&minute, &host, &max_usage, &min_usage)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 				// q.Logger.Info(minute, host, max_usage, min_usage)
// 			}

// 		}
// 	}

// 	return &results.Result{QueryTimes: queryTimes, Error: nil}
// }
