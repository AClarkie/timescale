package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/AClarkie/timescale/pkg/csv"
	"github.com/AClarkie/timescale/pkg/querier"
	"github.com/blendle/zapdriver"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	DebugHost      string = "0.0.0.0:8090"
	queryParams    string
	goroutineCount int
	dbHost         string
	dbName         string
	dbUser         string
	dbPassword     string
	dbSSLMode      string
)

func init() {
	flag.StringVar(&queryParams, "queryParams", "query_params.csv", "Path to a query_params.csv")
	flag.IntVar(&goroutineCount, "goroutineCount", 10, "The number of goroutines to create")
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
	{
		if l, err := zapdriver.NewProduction(); err != nil {
			return errors.Wrap(err, "creating logger")
		} else {
			logger = l.Sugar()
		}
	}
	// Flush logs at the end of the applications lifetime
	defer logger.Sync()

	// Set config file location for local testing
	queryData, err := csv.ProcessCsv(queryParams)
	if err != nil {
		return errors.Wrap(err, "unable to process input csv")
	}

	logger.Infow("Application starting")
	defer logger.Info("Application finished")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Prometheus Setup
	go func() {
		logger.Infow("metrics listener starting", "addr", DebugHost)
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(DebugHost, http.DefaultServeMux)
		logger.Fatalf("metrics listener closed", "err", err)
	}()

	// Database connection
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s dbname=%s sslmode=%s user=%s password=%s", dbHost, dbName, dbSSLMode, dbUser, dbPassword))
	defer db.Close()
	if err != nil {
		return errors.Wrap(err, "failed to setup database connection")
	}

	var wg sync.WaitGroup

	// iterate over 10 goroutines assigning hostnames to each
	for i, hostname := range allocateHostToExecutor(queryData, goroutineCount) {
		wg.Add(1)

		go func(hostname []csv.Hostname, db *sql.DB) {
			defer wg.Done()
			querier, err := querier.NewQuerier(hostname, logger.Named(fmt.Sprintf("executor %d", i)), db)
			if err != nil {
				logger.Fatal("Failed to create querier")
			}
			querier.Start()
		}(hostname, db)
	}
	wg.Wait()

	return nil
}

func allocateHostToExecutor(queryData map[string]csv.Hostname, goroutineCount int) [][]csv.Hostname {
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
