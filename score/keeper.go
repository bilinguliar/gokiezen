package score

import "github.com/mediocregopher/radix.v2/redis"

// Listing of Redis commands that we need to work with sets.
const (
	countries = "ALL_COUNTRIES"
	parties   = "ALL_PARTIES"

	redisGet       = "GET"
	redisIncr      = "INCR"
	redisSAdd      = "SADD"
	redisSRem      = "SREM"
	redisSMembers  = "SMEMBERS"
	redisSIsMember = "SISMEMBER"
)

// Connection pool used to make requests to Redis. Its size is provided by Configuration during application init.
type ConnectionPool interface {
	Cmd(cmd string, args ...interface{}) *redis.Resp
}

// Keeper is an implemetation of ScoreKeeper that uses Redis.
type Keeper struct {
	pool ConnectionPool
}

// NewKeeper returns pointer to created Keeper instance initialized with Redis pool.
func NewKeeper(p ConnectionPool) *Keeper {
	return &Keeper{pool: p}
}

// Get returns current score for given key.
func (d Keeper) Get(key string) (int, error) {
	return d.pool.Cmd(redisGet, key).Int()
}

// AddPoint increments counter by one for a given key.
func (d Keeper) AddPoint(key string) error {
	_, err := d.pool.Cmd(redisIncr, key).Int()
	return err
}

// AddCountry will create country record in set of all countries.
func (d Keeper) AddCountry(c string) error {
	return d.sadd(countries, c)
}

// AddCandidate adds the one to current voting.
func (d Keeper) AddCandidate(p string) error {
	return d.sadd(parties, p)
}

// RemoveCandidate deletes single candidate with specified name.
func (d Keeper) RemoveCandidate(p string) error {
	return d.srem(parties, p)
}

// GetAllCandidates returns all candidates currently taking part in voting.
func (d Keeper) GetAllCandidates() ([]string, error) {
	return d.smembers(parties)
}

// GetAllCountries returnes all countries that were participating during voting.
func (d Keeper) GetAllCountries() ([]string, error) {
	return d.smembers(countries)
}

func (d Keeper) sadd(set, name string) error {
	_, err := d.pool.Cmd(redisSAdd, set, name).Int()
	return err
}

func (d Keeper) srem(set, name string) error {
	_, err := d.pool.Cmd(redisSRem, set, name).Int()
	return err
}

// smembers WILL SLOWDOWN your SERVER if used with large sets.
func (d Keeper) smembers(set string) ([]string, error) {
	response, err := d.pool.Cmd(redisSMembers, set).List()
	return response, err
}
