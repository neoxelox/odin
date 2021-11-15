package database

import (
	"regexp"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/neoxelox/odin/internal"
)

const PGERRCODE_PATTERN = `\(SQLSTATE (.*)\)`

var (
	ErrNoRows             = internal.NewError("No rows in result set")
	ErrIntegrityViolation = internal.NewError("Integrity constraint violation")

	pgerrcodeRegex = regexp.MustCompile(PGERRCODE_PATTERN)
)

func internalError(err error) error {
	if err == nil {
		return nil
	}

	if code := pgerrcodeRegex.FindStringSubmatch(err.Error()); len(code) == 2 {
		switch code[1] {
		case pgerrcode.IntegrityConstraintViolation, pgerrcode.RestrictViolation, pgerrcode.NotNullViolation,
			pgerrcode.ForeignKeyViolation, pgerrcode.UniqueViolation, pgerrcode.CheckViolation,
			pgerrcode.ExclusionViolation:
			return ErrIntegrityViolation().WrapWithDepth(2, err)
		}
	}

	switch err.Error() {
	case pgx.ErrNoRows.Error():
		return ErrNoRows().WrapWithDepth(2, err)
	default:
		return err
	}
}
