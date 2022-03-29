package querier

import (
	"regexp"
	"testing"
	"time"

	"github.com/AClarkie/timescale/pkg/csv"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
)

func TestExecute(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Logging setup
	var logger *zap.SugaredLogger
	if l, err := zapdriver.NewProduction(); err != nil {
		t.Fatalf("an error '%s' was not expected when creating the logger", err)
	} else {
		logger = l.Sugar()
	}

	// Flush logs at the end of test
	defer logger.Sync()

	pTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
	hostnames := []csv.Hostname{{Name: "1", Times: []csv.Time{{StartTime: pTime, EndTime: pTime}}}}

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(hostnames[0].Name, hostnames[0].Times[0].StartTime, hostnames[0].Times[0].EndTime).WillReturnResult(sqlmock.NewResult(1, 1))

	result := Execute(hostnames, logger.Named("testing"), db)
	if result.Error != nil {
		t.Fatalf("an error '%s' was not expected when querying the database", err)
	}
}
