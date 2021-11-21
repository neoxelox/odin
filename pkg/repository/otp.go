package repository

import (
	"context"
	"fmt"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg/model"
)

const OTP_TABLE = "otp"

var ErrOTPGeneric = internal.NewError("OTP query failed")

type OTPRepository struct {
	class.Repository
}

func NewOTPRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *OTPRepository {
	return &OTPRepository{
		Repository: *class.NewRepository(configuration, logger, database),
	}
}

func (self *OTPRepository) Create(ctx context.Context, otp model.OTP) (*model.OTP, error) {
	var o model.OTP

	query := fmt.Sprintf(`INSERT INTO "%s"
						  ("id", "asset", "type", "code", "attempts", "expires_at")
						  VALUES ($1, $2, $3, $4, $5, $6)
						  RETURNING *;`, OTP_TABLE)

	err := self.Database.Query(
		ctx, query, otp.ID, otp.Asset, otp.Type, otp.Code, otp.Attempts, otp.ExpiresAt).Scan(&o)
	if err != nil {
		return nil, ErrOTPGeneric().Wrap(err)
	}

	return &o, nil
}

func (self *OTPRepository) GetByID(ctx context.Context, id string) (*model.OTP, error) {
	var o model.OTP

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "id" = $1;`, OTP_TABLE)

	err := self.Database.Query(ctx, query, id).Scan(&o)
	switch {
	case err == nil:
		return &o, nil
	case database.ErrNoRows().Is(err):
		return nil, nil
	default:
		return nil, ErrOTPGeneric().Wrap(err)
	}
}

func (self *OTPRepository) GetByAsset(ctx context.Context, asset string) (*model.OTP, error) {
	var o model.OTP

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "asset" = $1;`, OTP_TABLE)

	err := self.Database.Query(ctx, query, asset).Scan(&o)
	switch {
	case err == nil:
		return &o, nil
	case database.ErrNoRows().Is(err):
		return nil, nil
	default:
		return nil, ErrOTPGeneric().Wrap(err)
	}
}

func (self *OTPRepository) UpdateAttempts(ctx context.Context, id string, attempts int) error {
	query := fmt.Sprintf(`UPDATE "%s"
						  SET "attempts" = $1
						  WHERE "id" = $2;`, OTP_TABLE)

	affected, err := self.Database.Exec(ctx, query, attempts, id)
	if err != nil {
		return ErrOTPGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrOTPGeneric()
	}

	return nil
}

func (self *OTPRepository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM "%s"
						  WHERE "id" = $1;`, OTP_TABLE)

	affected, err := self.Database.Exec(ctx, query, id)
	if err != nil {
		return ErrOTPGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrOTPGeneric()
	}

	return nil
}
