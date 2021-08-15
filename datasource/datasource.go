package datasource

import (
	"github.com/rs/zerolog/log"
	"go-there/cache"
	"go-there/data"
	"go-there/database"
)

// DataSource represents the source of the user data (database+cache). It abstracts the caching process. Currently,
// only the path operations are cached
type DataSource struct {
	*database.DataBase
	*cache.Cache
}

// Init initializes a datasource from a *database.DataBase and *cache.Cache. The cache can be nil.
func Init(db *database.DataBase, cache *cache.Cache) *DataSource {
	return &DataSource{
		DataBase: db,
		Cache:    cache,
	}
}

// SelectUser fetches an complete user by his username in the database. Returns a data.ErrSql if it fails.
func (ds *DataSource) SelectUser(username string) (data.UserInfo, error) {
	return ds.DataBase.SelectUser(username)
}

// SelectAllUsers fetches the complete list of all users. Returns a data.ErrSql if it fails.
func (ds *DataSource) SelectAllUsers() ([]data.UserInfo, error) {
	return ds.DataBase.SelectAllUsers()
}

// SelectUserLogin fetches the id,username,is_admin,password_hash of a user by his username in the database. Returns a
// data.ErrSql if it fails.
func (ds *DataSource) SelectUserLogin(username string) (data.User, error) {
	return ds.DataBase.SelectUserLogin(username)
}

// SelectApiKeyHashByUser fetches a full API key hash from the database by a username. Returns a data.ErrSql if it
// fails.
func (ds *DataSource) SelectApiKeyHashByUser(username string) ([]byte, error) {
	return ds.DataBase.SelectApiKeyHashByUser(username)
}

// SelectUserLoginByApiKeyHash fetches the id,username,is_admin,api_key_hash of a user, by his API key hash.
func (ds *DataSource) SelectUserLoginByApiKeyHash(apiKeyHash string) (data.User, error) {
	return ds.DataBase.SelectUserLoginByApiKeyHash(apiKeyHash)
}

// InsertUser tries to add a new user to the database. If a user with the same name or API key hash exists,
// data.ErrSqlDuplicateRow is returned.
func (ds *DataSource) InsertUser(user data.User) error {
	return ds.DataBase.InsertUser(user)
}

// UpdateUserPassword updates an user's password in the database. Returns a data.ErrSql if it fails.
func (ds *DataSource) UpdateUserPassword(user data.User) error {
	return ds.DataBase.UpdateUserPassword(user)
}

// UpdateUserApiKey updates an user's API key in the database. Returns a data.ErrSql if it fails.
func (ds *DataSource) UpdateUserApiKey(user data.User) error {
	return ds.DataBase.UpdateUserApiKey(user)
}

// DeleteUser deletes a user in the database by his username. Returns a data.ErrSql if it fails.
// Logs a warning if a cache related error happens.
func (ds *DataSource) DeleteUser(username string) error {
	ui, err := ds.DataBase.SelectUser(username)

	if err != nil {
		return err
	}

	paths := make([]string, len(ui.Paths))

	for i := range ui.Paths {
		paths[i] = ui.Paths[i].Path
	}

	err = ds.Cache.DeleteTargets(paths)

	if err != nil {
		log.Warn().Err(err).Msg("error removing user targets from cache")
	}

	return ds.DataBase.DeleteUser(username)
}

// GetTarget tries to get a target from the cache, then from the database on a miss. Returns a data.ErrSqlNoRow if the
// target doesn't exist or data.ErrSql if it fails. The target is immediately added to the cache on a miss.
// Logs a warning if a cache related error happens.
func (ds *DataSource) GetTarget(path string) (string, error) {
	t, err := ds.Cache.GetTarget(path)

	if err != nil {
		log.Warn().Err(err).Msg("error getting target in cache")
	}

	if t != "" {
		return t, nil
	}

	t, err = ds.DataBase.GetTarget(path)

	if err != nil {
		return "", err
	}

	// On cache miss
	err = ds.Cache.AddTarget(data.Path{
		Path:   path,
		Target: t,
	})

	if err != nil {
		log.Warn().Err(err).Msg("error inserting path in cache")
	}

	return t, nil
}

// InsertPath adds a data.Path to the cache, then to the database. Returns a data.ErrSqlDuplicateRow if the path already
// exists or data.ErrSql if the operation fails.
// Logs a warning if a cache related error happens.
func (ds *DataSource) InsertPath(path data.Path) error {
	err := ds.Cache.AddTarget(path)

	if err != nil {
		log.Warn().Err(err).Msg("error inserting path in cache")
	}

	return ds.DataBase.InsertPath(path)
}

// DeletePath removes a data.Path from the cache, then deletes it in the database. Logs a warning if the cache returns
// an error, returns a data.ErrSql if the operation fails.
// Logs a warning if a cache related error happens.
func (ds *DataSource) DeletePath(path data.Path) error {
	err := ds.Cache.DeleteTargets([]string{path.Path})

	if err != nil {
		log.Warn().Err(err).Msg("error deleting path in cache")
	}

	return ds.DataBase.DeletePath(path)
}
