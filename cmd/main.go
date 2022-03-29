package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/AClarkie/timescale/pkg/csv"
	"github.com/AClarkie/timescale/pkg/querier"
	"github.com/AClarkie/timescale/pkg/results"
	"github.com/blendle/zapdriver"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	queryParams    string
	verbose        bool
	goroutineCount int
	dbHost         string
	dbName         string
	dbUser         string
	dbPassword     string
	dbSSLMode      string
)

func init() {
	flag.StringVar(&queryParams, "queryParams", "query_params.csv", "Path to the input csv")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.IntVar(&goroutineCount, "goroutineCount", 2, "The number of goroutines to use")
	flag.StringVar(&dbHost, "dbHost", "localhost", "The database host")
	flag.StringVar(&dbName, "dbName", "homework", "The database name")
	flag.StringVar(&dbUser, "dbUser", "postgres", "The database user")
	flag.StringVar(&dbPassword, "dbPassword", "password", "The database password")
	flag.StringVar(&dbSSLMode, "dbSSLMode", "disable", "The database sslmode")
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error : %s", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()

	// Logging setup
	var logger *zap.SugaredLogger
	if verbose {
		if l, err := zapdriver.NewDevelopment(); err != nil {
			return errors.Wrap(err, "creating verbose logger")
		} else {
			logger = l.Sugar()
		}
	} else {
		if l, err := zapdriver.NewProduction(); err != nil {
			return errors.Wrap(err, "creating production logger")
		} else {
			logger = l.Sugar()
		}
	}

	// Flush logs at the end of the applications lifetime
	defer logger.Sync()

	logger.Debug("Application starting")
	defer logger.Debug("Application finished")

	// Set config file location for local testing
	queryData, err := csv.ProcessCSV(queryParams)
	if err != nil {
		return errors.Wrap(err, "unable to process input csv")
	}

	// Database connection
	logger.Debug("Establishing connection to the database")
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s dbname=%s sslmode=%s user=%s password=%s", dbHost, dbName, dbSSLMode, dbUser, dbPassword))
	if err != nil {
		return errors.Wrap(err, "failed to setup database connection")
	}
	defer db.Close()

	// Setup waitgroup and results channel
	var wg sync.WaitGroup
	resultsChan := make(chan *results.Result)

	logger.Debug("Assigning hosts to groups and starting goroutines")
	for i, hostname := range assignHostsToGroups(queryData, goroutineCount) {
		wg.Add(1)

		go func(hostname []csv.Hostname, db *sql.DB, routine int) {
			defer wg.Done()
			resultsChan <- querier.Execute(hostname, logger.Named(fmt.Sprintf("querier %d", routine)), db)
		}(hostname, db, i)
	}

	// Launch a goroutine to monitor when all the work is done.
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Process and display results
	results.ProcessAndDisplayResults(resultsChan)

	return nil
}

// Assigns hosts to groups based on the number of go routines desired such that if
// there are 10 routines and 20 hosts then each processing group will contain 2 hosts
func assignHostsToGroups(queryData map[string]csv.Hostname, goroutineCount int) [][]csv.Hostname {
	groups := make([][]csv.Hostname, goroutineCount)

	i := 0
	for _, hostname := range queryData {
		if i == goroutineCount {
			i = 0
		}
		groups[i] = append(groups[i], hostname)
		i++
	}
	return groups
}
