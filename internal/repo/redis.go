package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dusk-chancellor/dc-sso/internal/models"
	"github.com/redis/go-redis/v9"
)

// redis caching methods implementation

// cached data expires after a day of no use
var expirationTime = 24 * time.Hour

// key type: `user:%s:%v`, where:
// `user` stands for what is being cached (user data here)
// `%s` key field (type) -> {id / email / username}
// `%v` key value
// e.g `user:Email:example@example.com`

// retrieves cached user model by its key field
// if data wasn't cached in redis
// calls db method
func (r *Rdb) GetUserBy(ctx context.Context, field Field, val any) (models.User, error) {
	key := fmt.Sprintf("user:%s:%v", field, val)

	res, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil { // if not cached
		user, err := r.db.GetUserBy(ctx, field, val)
		if err != nil {
			return models.User{}, err
		}
		// user model stored as []byte
		userData, err := json.Marshal(user)
		if err != nil {
			return models.User{}, err
		}

		if err := r.rdb.Set(ctx, key, userData, expirationTime).Err(); err != nil {
			return models.User{}, err
		}

		return user, nil
	}

	if err != nil {
		return models.User{}, err
	}

	var user models.User

	if err := json.Unmarshal([]byte(res), &user); err != nil {
		return models.User{}, err
	}

	return user, nil
}

// updates user info both in db & redis
func (r *Rdb) UpdateByID(ctx context.Context, id string, field Field, val any) error {
	key := fmt.Sprintf("user:%s:%v", ID, id)

	// update in db
	if err := r.db.UpdateByID(ctx, id, field, val); err != nil {
		return err
	}
	// get updated user info
	user, err := r.db.GetUserBy(ctx, ID, id)
	if err != nil {
		return err
	}
	// transform into []byte type
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	// set updated data in cache
	if err := r.rdb.Set(ctx, key, userData, expirationTime).Err(); err != nil {
		return err
	}

	return nil
}

// holds the same logic as `UpdateByID` but requires `email` instead
func (r *Rdb) UpdateRole(ctx context.Context, email, wantsRole string) error {
	key := fmt.Sprintf("user:%s:%v", Email, email)

	if err := r.db.UpdateRole(ctx, email, wantsRole); err != nil {
		return err
	}

	user, err := r.db.GetUserBy(ctx, Email, email)
	if err != nil {
		return err
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := r.rdb.Set(ctx, key, userData, expirationTime).Err(); err != nil {
		return err
	}

	return nil
}
