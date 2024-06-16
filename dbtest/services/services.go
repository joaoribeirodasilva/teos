package services

import (
	"github.com/joaoribeirodasilva/teos/dbtest/database"
	"github.com/joaoribeirodasilva/teos/dbtest/memdb"
)

// TODO: PENSAR MELHOR ser serviços vão estar em classes, talvez sim
var (
	services Services
)

type Services struct {
	memDbs   *map[string]*memdb.RedisDB
	mongoDbs *map[string]*database.Db
	// servers
	// user
	// http query data
	// configuration
	// conf
}

func SetMongo(svcs *map[string]*database.Db) {

	services.mongoDbs = svcs
}

func SetMemDbs(svcs *map[string]*memdb.RedisDB) {

	services.memDbs = svcs
}

func SetServers() {
	// TODO:
}

func GetMongo(name string) *database.Db {

	//db := memDbs
	return nil
}
