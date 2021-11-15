package cache

import (
	"github.com/go-redis/cache/v8"

	"github.com/neoxelox/odin/internal"
)

var (
	ErrMiss = internal.NewError("Key not present")
)

func internalError(err error) error {
	if err == nil {
		return nil
	}

	switch err {
	case cache.ErrCacheMiss:
		return ErrMiss().WrapWithDepth(2, err)
	default:
		return err
	}
}
