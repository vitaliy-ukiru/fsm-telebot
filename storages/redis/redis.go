// Package redis contains redis storage.
package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/redis/go-redis/v9"
	"github.com/vitaliy-ukiru/fsm-telebot"
)

type keyType string

const (
	stateKey     keyType = "state"
	stateDataKey keyType = "data"
)

type Storage struct {
	rds  *redis.Client
	pref StorageSettings
}

type StorageSettings struct {
	// Prefix for records in Redis.
	// Default is "fsm".
	Prefix string

	// TTL for state.
	// Default is 0 (no ttl).
	TTLState time.Duration

	// TTL for state data.
	// Default is 0 (no ttl).
	TTLData time.Duration

	// Batch size for reset data.
	// Default is 0 (no batching).
	ResetDataBatchSize int64
}

// NewStorage returns new redis storage.
func NewStorage(client *redis.Client, pref StorageSettings) *Storage {
	if pref.Prefix == "" {
		pref.Prefix = "fsm"
	}

	return &Storage{
		rds:  client,
		pref: pref,
	}
}

func (s *Storage) GetState(chatId, userId int64) (fsm.State, error) {
	val, err := s.rds.Get(context.TODO(), s.generateKey(chatId, userId, stateKey)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fsm.DefaultState, nil
		}
		return fsm.DefaultState, wrapError(err, "get")
	}

	return fsm.State(val), nil
}

func (s *Storage) SetState(chatId, userId int64, state fsm.State) error {
	err := s.rds.Set(
		context.TODO(),
		s.generateKey(chatId, userId, stateKey),
		string(state),
		s.pref.TTLState,
	).Err()
	return wrapError(err, "set")
}

func (s *Storage) ResetState(chatId, userId int64, withData bool) error {
	if err := s.SetState(chatId, userId, fsm.DefaultState); err != nil {
		return wrapError(err, "set state to default")
	}

	if withData {
		if err := s.resetData(chatId, userId); err != nil {
			return wrapError(err, "reset data")
		}
	}
	return nil
}

func (s *Storage) resetData(chatId, userId int64) error {
	var cursor uint64
	var keys []string

	redisKey := s.generateKey(chatId, userId, stateDataKey, "*")

	for {
		var err error
		keys, cursor, err = s.rds.Scan(
			context.TODO(),
			cursor,
			redisKey,
			s.pref.ResetDataBatchSize,
		).Result()
		if err != nil {
			return errors.Wrap(err, "scan")
		}

		if len(keys) > 0 {
			if err := s.rds.Del(context.TODO(), keys...).Err(); err != nil {
				return errors.Wrap(err, "delete keys")
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (s *Storage) UpdateData(chatId, userId int64, key string, data interface{}) error {
	ctx := context.TODO()
	redisKey := s.generateKey(chatId, userId, stateDataKey, key)

	if data == nil {
		err := s.rds.Del(ctx, redisKey).Err()
		return wrapError(err, "delete data")
	}

	encodedData, err := s.encode(data)
	if err != nil {
		return wrapError(err, "encode data")
	}

	err = s.rds.Set(
		ctx,
		redisKey,
		encodedData,
		s.pref.TTLData,
	).Err()

	return wrapError(err, "set data")
}

func (s *Storage) GetData(chatId, userId int64, key string, to interface{}) error {
	dataBytes, err := s.rds.Get(context.TODO(), s.generateKey(chatId, userId, stateDataKey, key)).Bytes()
	if errors.Is(err, redis.Nil) {
		return fsm.ErrNotFound
	}
	if err != nil {
		return wrapError(err, "get data")
	}

	if err := s.decode(dataBytes, to); err != nil {
		return wrapError(err, "decode data")
	}
	return nil
}

func (s *Storage) Close() error {
	return s.rds.Close()
}

func (s *Storage) generateKey(chat, user int64, keyType keyType, keys ...string) string {
	res := fmt.Sprintf("%s:%d:%d:%s", s.pref.Prefix, chat, user, keyType)
	if len(keys) > 0 {
		res += ":" + strings.Join(keys, ":")
	}
	return res
}

func (s *Storage) encode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *Storage) decode(data []byte, to interface{}) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(to)
}

func wrapError(err error, msg string) error {
	return errors.Wrapf(err, "fsm-telebot/storage/redis: %s", msg)
}
