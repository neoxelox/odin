package repository

import (
	"context"
	"fmt"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/rs/xid"
)

const OTP_TABLE = "otp"

var ErrOTPGeneric = internal.NewError("OTP query failed")

type OTPRepository struct {
	class.Repository
}

func NewOTPRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *OTPRepository {
	return &OTPRepository{
		Repository: *class.NewRepository(OTP_TABLE, configuration, logger, database),
	}
}

func (self *OTPRepository) Transaction(ctx context.Context, fn func(*OTPRepository) error) error {
	return self.Database.Transaction(ctx, func(db *database.Database) error {
		return fn(&OTPRepository{
			Repository: *class.NewRepository(self.Table, self.Configuration, self.Logger, *db),
		})
	})
}

func (self *OTPRepository) Create(ctx context.Context, otp model.OTP) (*model.OTP, error) {
	var o model.OTP

	query := fmt.Sprintf(`INSERT INTO "%s" ("id", "asset", "code", "attempts", "created_at", "updated_at", "deleted_at")
						  VALUES ($1, $2, $3, $4, $5, $6, $7)
						  RETURNING *;`, OTP_TABLE)

	err := self.Database.Query(ctx, query, otp.ID, otp.Asset, otp.Code, otp.Attempts, otp.CreatedAt, otp.UpdatedAt, otp.DeletedAt).Scan(&o)
	if err != nil {
		return nil, ErrOTPGeneric().Wrap(err)
	}

	return &o, nil
}

func (self *OTPRepository) Get(ctx context.Context, id xid.ID) (*model.OTP, error) {
	var o model.OTP

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE id = $1;`, OTP_TABLE)

	err := self.Database.Query(ctx, query, id).Scan(&o)
	if err != nil {
		return nil, ErrOTPGeneric().Wrap(err)
	}

	return &o, nil
}
