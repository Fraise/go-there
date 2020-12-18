package datasource

import (
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
func (ds *DataSource) SelectUser(username string) (data.User, error) {
	return ds.DataBase.SelectUser(username)
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

// SelectUserLoginByApiKeySalt fetches the id,username,is_admin,api_key_hash of a user, by his API key salt.
func (ds *DataSource) SelectUserLoginByApiKeySalt(apiKeySalt string) (data.User, error) {
	return ds.DataBase.SelectUserLoginByApiKeySalt(apiKeySalt)
}

// SelectApiKeyHashBySalt fetches the full API key hash by its salt. Returns a data.ErrSql if it fails.
func (ds *DataSource) SelectApiKeyHashBySalt(apiKeySalt string) ([]byte, error) {
	return ds.DataBase.SelectApiKeyHashBySalt(apiKeySalt)
}

// InsertUser tries to add a new user to the database. If a user with the same name or API key salt exists,
// data.ErrSqlDuplicateRow is returned.
func (ds *DataSource) InsertUser(user data.User) error {
	return ds.DataBase.InsertUser(user)
}

// UpdateUserPassword updates an user's password in the database. Returns a data.ErrSql if it fails.
func (ds *DataSource) UpdateUserPassword(user data.User) error {
	return ds.DataBase.UpdateUserPassword(user)
}

// UpdateUserPassword updates an user's API key in the database. Returns a data.ErrSql if it fails.
func (ds *DataSource) UpdateUserApiKey(user data.User) error {
	return ds.DataBase.UpdateUserApiKey(user)
}

// DeleteUser deletes a user in the database by his username. Returns a data.ErrSql if it fails.
func (ds *DataSource) DeleteUser(username string) error {
	return ds.DataBase.DeleteUser(username)
}

// GetTarget gets a target in the database from a path. Returns a data.ErrSqlNoRow if the target doesn't exist or
// data.ErrSql if it fails.
func (ds *DataSource) GetTarget(path string) (string, error) {
	return ds.DataBase.GetTarget(path)
}

// InsertPath adds a data.Path to the database. Returns a data.ErrSqlDuplicateRow if the path already exists or
// data.ErrSql if it fails.
func (ds *DataSource) InsertPath(path data.Path) error {
	return ds.DataBase.InsertPath(path)
}

// InsertPath deletes a data.Path in the database. Returns a data.ErrSql if it fails.
func (ds *DataSource) DeletePath(path data.Path) error {
	return ds.DataBase.DeletePath(path)
}
