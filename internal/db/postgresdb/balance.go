package postgresdb

import (
	"database/sql"
	"fmt"

	"github.com/siraj18/avito-test/internal/models"
)

var ErrorUserNotFound = fmt.Errorf("user not found")

func (rep *SqlRepository) GetBalance(uid string) (*models.User, error) {
	var user models.User

	if err := rep.db.Get(&user, getUserBalanceSql, uid); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorUserNotFound
		}
		return nil, fmt.Errorf("error when get user balance: %w", err)
	}

	return &user, nil

}

func (rep *SqlRepository) ChangeBalance(uid string, money float64) {

}

func (rep *SqlRepository) TransferBalance(fromUid string, toUid string, money float64) {

}
