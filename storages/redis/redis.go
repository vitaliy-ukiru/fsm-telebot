// Package redis contains redis storage.
package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vitaliy-ukiru/fsm-telebot"
)

type keyType string

const (
	stateKey     keyType = "state"
	stateDataKey keyType = "data"
)

type Storage struct {
	rds                *redis.Client
	prefix             string
	ttlState           time.Duration
	ttlData            time.Duration
	resetDataBatchSize int64
}

type StorageSettings struct {
	// Prefix for records in Redis.
	// Default is "fsm".
	prefix string

	// TTL for state.
	// Default is 0 (no ttl).
	ttlState time.Duration

	// TTL for state data.
	// Default is 0 (no ttl).
	ttlData time.Duration

	// Batch size for reset data.
	// Default is 0 (no batching).
	resetDataBatchSize int64
}

// NewStorage returns new redis storage.
func NewStorage(client *redis.Client, pref StorageSettings) fsm.Storage {
	if pref.prefix == "" {
		pref.prefix = "fsm"
	}

	return &Storage{
		rds:                client,
		prefix:             pref.prefix,
		ttlState:           pref.ttlState,
		ttlData:            pref.ttlData,
		resetDataBatchSize: pref.resetDataBatchSize,
	}
}

func (s *Storage) GetState(chatId, userId int64) (fsm.State, error) {
	val, err := s.rds.Get(context.TODO(), s.generateKey(chatId, userId, stateKey)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fsm.DefaultState, nil
		}
		return fsm.DefaultState, err
	}

	return fsm.State(val), nil
}

func (s *Storage) SetState(chatId, userId int64, state fsm.State) error {
	return s.rds.Set(context.TODO(), s.generateKey(chatId, userId, stateKey), string(state), s.ttlState).Err()
}

func (s *Storage) ResetState(chatId, userId int64, withData bool) error {
	if err := s.SetState(chatId, userId, fsm.DefaultState); err != nil {
		return err
	}

	if withData {
		s.resetData(chatId, userId)
	}
	return nil
}

func (s *Storage) resetData(chatId, userId int64) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = s.rds.Scan(
			context.TODO(),
			cursor,
			s.generateKey(chatId, userId, stateDataKey, "*"),
			s.resetDataBatchSize,
		).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := s.rds.Del(context.TODO(), keys...).Err(); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (s *Storage) UpdateData(chatId, userId int64, key string, data interface{}) error {
	if data == nil {
		return s.rds.Del(context.TODO(), s.generateKey(chatId, userId, stateDataKey, key)).Err()
	}

	encodedData, err := s.encode(data)
	if err != nil {
		return err
	}
	return s.rds.Set(context.TODO(), s.generateKey(chatId, userId, stateDataKey, key), encodedData, s.ttlData).Err()
}

func (s *Storage) GetData(chatId, userId int64, key string, to interface{}) error {
	dataBytes, err := s.rds.Get(context.TODO(), s.generateKey(chatId, userId, stateDataKey, key)).Bytes()
	if errors.Is(err, redis.Nil) {
		return fsm.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("get data: %w", err)
	}

	return s.decode(dataBytes, to)
}

func (s *Storage) Close() error {
	return s.rds.Close()
}

func (s *Storage) generateKey(chat, user int64, keyType keyType, keys ...string) string {
	res := fmt.Sprintf("%s:%d:%d:%s", s.prefix, chat, user, keyType)
	if keys != nil {
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
