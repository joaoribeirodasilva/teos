package memdb

type MemDBObptions struct {
	Addr     string
	Db       int
	Username string
	Password string
}

var (
	defaultOptions = MemDBObptions{
		Addr:     "localhost:6379",
		Db:       0,
		Username: "",
		Password: "",
	}
	dbs map[string]*RedisDB
)

func New(name string, options *MemDBObptions) *RedisDB {

	if dbs == nil {
		dbs = make(map[string]*RedisDB)
	}
	if options == nil {
		options = &defaultOptions
	}

	r := &RedisDB{
		addr:     options.Addr,
		db:       options.Db,
		username: options.Username,
		password: options.Password,
	}

	dbs[name] = r

	return r

}

func GetDatabase(name string) *RedisDB {
	db, ok := dbs[name]
	if !ok {
		return nil
	}
	return db
}

func GetAvailable() *[]string {

	if len(dbs) == 0 {
		return nil
	}

	keys := make([]string, 0)
	for k := range dbs {
		keys = append(keys, k)
	}

	return &keys
}
