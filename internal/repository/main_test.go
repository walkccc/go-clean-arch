package repository

import (
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/walkccc/go-clean-arch/internal/util"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	testDB := NewDB(config.DBDriver, config.DBSource)
	testQueries = New(testDB)
	os.Exit(m.Run())
}
