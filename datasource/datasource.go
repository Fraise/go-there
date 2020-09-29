package datasource

import (
	"fmt"
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

func (ds *DataSource) SelectUser(username string) data.User {
	return data.User{}
}

func (ds *DataSource) InsertUser(user data.User) error {
	return nil
}

func (ds *DataSource) DeleteUser(username string) error {
	return nil
}
