package datasource

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"go-there/config"
	"go-there/data"
)

// DataSource represents the source of the applciation's data.
type DataSource struct {
	db *sqlx.DB
}

// Init initializes and tries to connect to the database defined in the configuration. If it cannot connect, an error is
// returned.
func Init(config *config.Configuration) (*DataSource, error) {
	var err error
	ds := new(DataSource)

	ds.db, err = sqlx.Connect(
		config.Database.Type,
		fmt.Sprintf(
			"%s:%s@%s(%s:%d)/%s",
			config.Database.User,
			config.Database.Password,
			config.Database.Protocol,
			config.Database.Address,
			config.Database.Port,
			config.Database.Name,
		),
	)

	if err != nil {
		return nil, err
	}

	return ds, nil
}

// SelectUser fetches an complete user by his username in the database.
func (ds *DataSource) SelectUser(username string) (data.User, error) {
	u := data.User{}
	err := ds.db.Get(&u, ds.db.Rebind("SELECT * FROM users WHERE username=?"), username)

	if err != nil {
		return data.User{}, err
	}

	return u, nil
}

// SelectUserLogin fetches the username,is_admin,password_hash as a user by his username in the database.
func (ds *DataSource) SelectUserLogin(username string) (data.User, error) {
	u := data.User{}
	err := ds.db.Get(&u, ds.db.Rebind("SELECT username,is_admin,password_hash FROM users WHERE username=?"), username)

	if err != nil {
		return data.User{}, err
	}

	return u, nil
}

// SelectApiKeyHashByUser fetches a full API key hash from the database by a username.
func (ds *DataSource) SelectApiKeyHashByUser(username string) ([]byte, error) {
	ak := make([]byte, 0)
	err := ds.db.Get(&ak, ds.db.Rebind("SELECT api_key_hash FROM users WHERE username=?"), username)

	if err != nil {
		return []byte{}, err
	}

	return ak, nil
}

// SelectUserLoginByApiKeySalt fetches the username,is_admin,api_key_hash of a user, by his API key salt.
func (ds *DataSource) SelectUserLoginByApiKeySalt(apiKeySalt string) (data.User, error) {
	u := data.User{}
	err := ds.db.Get(&u, ds.db.Rebind("SELECT username,is_admin,api_key_hash FROM users WHERE api_key_salt=?"), apiKeySalt)

	if err != nil {
		return data.User{}, err
	}

	return u, nil
}

// SelectApiKeyHashBySalt fetches the full API key hash by its salt.
func (ds *DataSource) SelectApiKeyHashBySalt(apiKeySalt string) ([]byte, error) {
	ak := make([]byte, 0)
	err := ds.db.Get(&ak, ds.db.Rebind("SELECT api_key_hash FROM users WHERE api_key_salt=?"), apiKeySalt)

	if err != nil {
		return []byte{}, err
	}

	return ak, nil
}

// InsertUser tries to add a new user to the database. If a user with the same name or API key salt exists,
// data.ErrSqlDuplicateRow is returned.
func (ds *DataSource) InsertUser(user data.User) error {
	_, err := ds.db.NamedExec(
		"INSERT INTO users (username,is_admin,password_hash,api_key_salt,api_key_hash) "+
			"VALUES (:username,:is_admin,:password_hash,:api_key_salt,:api_key_hash)", user)

	if err != nil {
		if e, ok := err.(*mysql.MySQLError); ok {
			// mysql duplicate row
			if e.Number == 1062 {
				return data.ErrSqlDuplicateRow
			}
		} else {
			return data.ErrSql
		}
	}

	return nil
}

// UpdatetUserPassword updates an user's password in the database.
func (ds *DataSource) UpdatetUserPassword(user data.User) error {
	_, err := ds.db.NamedExec("UPDATE users SET password_hash=:password_hash WHERE username=:username", user)

	return err
}

// UpdatetUserPassword updates an user's API key in the database.
func (ds *DataSource) UpdatetUserApiKey(user data.User) error {
	_, err := ds.db.NamedExec("UPDATE users SET api_key_hash=:api_key_hash,api_key_salt=:api_key_salt WHERE username=:username", user)

	return err
}

// DeleteUser deletes a user in the database by his userame.
func (ds *DataSource) DeleteUser(username string) error {
	_, err := ds.db.Exec(ds.db.Rebind("DELETE FROM users WHERE username=?"), username)

	if err != nil {
		return err
	}

	return nil
}

// GetTarget gets a target in the database from a path.
func (ds *DataSource) GetTarget(path string) (string, error) {
	t := ""
	err := ds.db.Get(&t, ds.db.Rebind("SELECT target FROM go WHERE path=?"), path)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", data.ErrSqlNoRow
		default:
			return "", err
		}
	}

	return t, nil
}

// InsertPath adds a data.Path to the database.
func (ds *DataSource) InsertPath(path data.Path) error {
	_, err := ds.db.NamedExec("INSERT INTO go (path,target,user,is_public) VALUES (:path,:target,:user,:is_public)", path)

	if err != nil {
		if e, ok := err.(*mysql.MySQLError); ok {
			// mysql duplicate row
			if e.Number == 1062 {
				return data.ErrSqlDuplicateRow
			}
		} else {
			return data.ErrSql
		}

	}

	return err
}

// InsertPath deletes a data.Path in the database.
func (ds *DataSource) DeletePath(path data.Path) error {
	_, err := ds.db.NamedExec("DELETE FROM go WHERE path=:path", path)

	return err
}
