package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"go-there/config"
	"go-there/data"
)

// DataBase represents the database containing the application's data.
type DataBase struct {
	db *sqlx.DB
}

// Init initializes and tries to connect to the database defined in the configuration. If it cannot connect, an error is
// returned.
func Init(config *config.Configuration) (*DataBase, error) {
	ds, err := connect(config, config.Database.Type)

	if err != nil {
		return nil, err
	}

	return ds, nil
}

// connect tries to connect to a database with the specified parameters. Returns data.ErrSql if it fails.
func connect(config *config.Configuration, dbType string) (*DataBase, error) {
	var err error
	ds := new(DataBase)

	switch dbType {
	case "mysql":
		ds.db, err = sqlx.Connect(
			dbType,
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
	case "postgres":
		// TODO
		log.Fatal().Err(errors.New("not implemented"))
		ds.db, err = sqlx.Connect(
			dbType,
			fmt.Sprintf(
				"user=%s password=%s host=%s port=%d name=%s sslmode=%s",
				config.Database.User,
				config.Database.Password,
				config.Database.Address,
				config.Database.Port,
				config.Database.Name,
				func() string {
					if config.Database.SslMode {
						return "enable"
					} else {
						return "disable"
					}
				}(),
			),
		)
	default:
		return nil, fmt.Errorf("%w : %s", data.ErrSql, "invalid sql type")
	}

	if err != nil {
		return ds, fmt.Errorf("%w : %s", data.ErrSql, err)
	}

	return ds, nil
}

// SelectUser fetches an complete user by his username in the database. Returns a data.ErrSql if it fails.
func (ds *DataBase) SelectUser(username string) (data.User, error) {
	u := data.User{}
	err := ds.db.Get(&u, ds.db.Rebind("SELECT * FROM users WHERE username=?"), username)

	if err != nil {
		return data.User{}, fmt.Errorf("%w : %s", data.ErrSql, err)
	}

	return u, nil
}

// SelectUserLogin fetches the id,username,is_admin,password_hash of a user by his username in the database. Returns a
// data.ErrSql if it fails.
func (ds *DataBase) SelectUserLogin(username string) (data.User, error) {
	u := data.User{}
	err := ds.db.Get(&u, ds.db.Rebind("SELECT id,username,is_admin,password_hash FROM users WHERE username=?"), username)

	if err != nil {
		return data.User{}, fmt.Errorf("%w : %s", data.ErrSql, err)
	}

	return u, nil
}

// SelectApiKeyHashByUser fetches a full API key hash from the database by a username. Returns a data.ErrSql if it
// fails.
func (ds *DataBase) SelectApiKeyHashByUser(username string) ([]byte, error) {
	ak := make([]byte, 0)
	err := ds.db.Get(&ak, ds.db.Rebind("SELECT api_key_hash FROM users WHERE username=?"), username)

	if err != nil {
		return []byte{}, fmt.Errorf("%w : %s", data.ErrSql, err)
	}

	return ak, nil
}

// SelectUserLoginByApiKeySalt fetches the id,username,is_admin,api_key_hash of a user, by his API key salt.
func (ds *DataBase) SelectUserLoginByApiKeySalt(apiKeySalt string) (data.User, error) {
	u := data.User{}
	err := ds.db.Get(&u, ds.db.Rebind("SELECT id,username,is_admin,api_key_hash FROM users WHERE api_key_salt=?"), apiKeySalt)

	if err != nil {
		return data.User{}, fmt.Errorf("%w : %s", data.ErrSql, err)
	}

	return u, nil
}

// SelectApiKeyHashBySalt fetches the full API key hash by its salt. Returns a data.ErrSql if it fails.
func (ds *DataBase) SelectApiKeyHashBySalt(apiKeySalt string) ([]byte, error) {
	ak := make([]byte, 0)
	err := ds.db.Get(&ak, ds.db.Rebind("SELECT api_key_hash FROM users WHERE api_key_salt=?"), apiKeySalt)

	if err != nil {
		return []byte{}, fmt.Errorf("%w : %s", data.ErrSql, err)
	}

	return ak, nil
}

// InsertUser tries to add a new user to the database. If a user with the same name or API key salt exists,
// data.ErrSqlDuplicateRow is returned.
func (ds *DataBase) InsertUser(user data.User) error {
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
			return fmt.Errorf("%w : %s", data.ErrSql, err)
		}
	}

	return nil
}

// UpdateUserPassword updates an user's password in the database. Returns a data.ErrSql if it fails.
func (ds *DataBase) UpdateUserPassword(user data.User) error {
	_, err := ds.db.NamedExec("UPDATE users SET password_hash=:password_hash WHERE username=:username", user)

	if err != nil {
		return fmt.Errorf("%w : %s", data.ErrSql, err)
	}

	return nil
}

// UpdateUserPassword updates an user's API key in the database. Returns a data.ErrSql if it fails.
func (ds *DataBase) UpdateUserApiKey(user data.User) error {
	_, err := ds.db.NamedExec("UPDATE users SET api_key_hash=:api_key_hash,api_key_salt=:api_key_salt WHERE username=:username", user)

	if err != nil {
		return fmt.Errorf("%w : %s", data.ErrSql, err)
	}

	return nil
}

// DeleteUser deletes a user in the database by his username. Returns a data.ErrSql if it fails.
func (ds *DataBase) DeleteUser(username string) error {
	_, err := ds.db.Exec(ds.db.Rebind("DELETE FROM users WHERE username=?"), username)

	if err != nil {
		return fmt.Errorf("%w : %s", data.ErrSql, err)
	}

	return nil
}

// GetTarget gets a target in the database from a path. Returns a data.ErrSqlNoRow if the target doesn't exist or
// data.ErrSql if it fails.
func (ds *DataBase) GetTarget(path string) (string, error) {
	t := ""
	err := ds.db.Get(&t, ds.db.Rebind("SELECT target FROM go WHERE path=?"), path)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", data.ErrSqlNoRow
		default:
			return "", fmt.Errorf("%w : %s", data.ErrSql, err)
		}
	}

	return t, nil
}

// InsertPath adds a data.Path to the database. Returns a data.ErrSqlDuplicateRow if the path already exists or
// data.ErrSql if it fails.
func (ds *DataBase) InsertPath(path data.Path) error {
	_, err := ds.db.NamedExec("INSERT INTO go (path,target,user_id) VALUES (:path,:target,:user_id)", path)

	if err != nil {
		if e, ok := err.(*mysql.MySQLError); ok {
			// mysql duplicate row
			if e.Number == 1062 {
				return data.ErrSqlDuplicateRow
			}
		} else {
			return fmt.Errorf("%w : %s", data.ErrSql, err)
		}
	}

	return nil
}

// InsertPath deletes a data.Path in the database. Returns a data.ErrSql if it fails.
func (ds *DataBase) DeletePath(path data.Path) error {
	_, err := ds.db.NamedExec("DELETE FROM go WHERE path=:path AND user_id=:user_id", path)

	if err != nil {
		return fmt.Errorf("%w : %s", data.ErrSql, err)
	}

	return nil
}
