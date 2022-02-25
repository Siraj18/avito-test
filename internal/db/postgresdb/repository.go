package postgresdb

import "github.com/jmoiron/sqlx"

type SqlRepository struct {
	db *sqlx.DB
}

func NewSqlRepository(db *sqlx.DB) (*SqlRepository, error) {
	rep := &SqlRepository{db}

	if err := rep.init(); err != nil {
		return nil, err
	}

	return rep, nil
}

func (rep *SqlRepository) init() error {
	_, err := rep.db.Exec(initSchema)

	if err != nil {
		return err
	}

	return nil
}

func (rep *SqlRepository) Close() {
	rep.db.Close()
}
