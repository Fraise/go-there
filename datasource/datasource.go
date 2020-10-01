package datasource

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"go-there/config"
	"go-there/data"
)

type DataSource struct {
	db *sqlx.DB
}

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

func (ds *DataSource) SelectUser(username string) (data.User, error) {
	u := data.User{}
	err := ds.db.Get(&u, ds.db.Rebind("SELECT * FROM users WHERE username=?"), username)

	if err != nil {
		return data.User{}, err
	}

	return u, nil
}

func (ds *DataSource) SelectUserPassword(username string) ([]byte, error) {
	p := make([]byte, 0)
	err := ds.db.Get(&p, ds.db.Rebind("SELECT password_salt,password_hash FROM users WHERE username=?"), username)

	if err != nil {
		return []byte{}, err
	}

	return p, nil
}

func (ds *DataSource) SelectUserApiKey(username string) ([]byte, error) {
	ak := make([]byte, 0)
	err := ds.db.Get(&ak, ds.db.Rebind("SELECT api_key_salt,api_key_hash FROM users WHERE username=?"), username)

	if err != nil {
		return []byte{}, err
	}

	return ak, nil
}

func (ds *DataSource) SelectApiKey(apiKey string) ([]byte, error) {
	ak := make([]byte, 0)
	err := ds.db.Get(&ak, ds.db.Rebind("SELECT api_key_hash FROM users WHERE api_key_hash=?"), apiKey)

	if err != nil {
		return []byte{}, err
	}

	return ak, nil
}

func (ds *DataSource) InsertUser(user data.User) error {
	return nil
}

func (ds *DataSource) DeleteUser(username string) error {
	return nil
}
