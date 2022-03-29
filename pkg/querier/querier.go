package querier

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/AClarkie/timescale/pkg/csv"
	"github.com/AClarkie/timescale/pkg/results"
	"go.uber.org/zap"
)

var (
	query = `SELECT
		time_bucket('1 minute', ts) AS minute,
		host,
		max(usage) AS max_usage,
		min(usage) AS min_usage
		FROM cpu_usage
		WHERE ts >= $2 AND ts < $3
		AND host = $1
		GROUP BY minute, host
		ORDER BY minute DESC;`
)

func Execute(hostnames []csv.Hostname, logger *zap.SugaredLogger, db *sql.DB) *results.Result {
	var queryTimes []time.Duration
	for _, hostname := range hostnames {
		logger.Debugf("Starting queries for hostname %s", hostname.Name)
		for _, timeRange := range hostname.Times {

			start := time.Now()
			_, err := db.Exec(query, hostname.Name, timeRange.StartTime, timeRange.EndTime)

			if err != nil {
				logger.Debugf("Query failed for hostname %s", hostname.Name)
				return &results.Result{QueryTimes: nil, Error: err}
			}
			queryTimes = append(queryTimes, time.Since(start))
		}
		logger.Debugf("Finished queries for hostname %s", hostname.Name)
	}

	logger.Debug("Execution complete")
	return &results.Result{QueryTimes: queryTimes, Error: nil}
}
