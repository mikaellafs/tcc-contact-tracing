package repositories

import (
	"errors"
	"strings"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

func parsePostgreSQLError(err error) error {
	if strings.Contains(err.Error(), "no rows") {
		err = ErrDuplicate
	}

	return err
}
