// Package redis contains redis storage.
package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	rds      *redis.Client
	prefix   string
	ttlState time.Duration
	ttlData  time.Duration
}

type StorageSettings struct {
	// Prefix for records in Redis.
	// Default is "fsm".
	prefix string

	// TTL for state.
	// Default is 0.
	ttlState time.Duration

	// TTL for state data.
	// Default is 0.
	ttlData time.Duration
}

// NewStorage returns new redis storage.
func NewStorage(client *redis.Client, pref StorageSettings) fsm.Storage {
	if pref.prefix == "" {
		pref.prefix = "fsm"
	}

	return &Storage{
		rds:      client,
		prefix:   pref.prefix,
		ttlState: pref.ttlState,
		ttlData:  pref.ttlData,
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
		return s.rds.Del(context.TODO(), s.generateKey(chatId, userId, stateDataKey)).Err()
	}
	return nil
}

func (s *Storage) UpdateData(chatId, userId int64, key string, data interface{}) error {
	stateDataBytes, err := s.rds.Get(context.TODO(), s.generateKey(chatId, userId, stateDataKey)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			stateDataBytes = []byte("{}")
		} else {
			return fmt.Errorf("get state data: %w", err)
		}
	}

	// unmarshal state data
	var stateData map[string]interface{}
	if err := json.Unmarshal(stateDataBytes, &stateData); err != nil {
		return fmt.Errorf("unmarshal state data: %w", err)
	}

	// set or delete data
	if data == nil {
		delete(stateData, key)
	} else {
		stateData[key] = data
	}

	// marshal state data
	stateDataBytes, err = json.Marshal(stateData)
	if err != nil {
		return fmt.Errorf("marshal state data: %w", err)
	}

	return s.rds.Set(context.TODO(), s.generateKey(chatId, userId, stateDataKey), stateDataBytes, s.ttlData).Err()
}

func (s *Storage) GetData(chatId, userId int64, key string) (interface{}, error) {
	stateDataBytes, err := s.rds.Get(context.TODO(), s.generateKey(chatId, userId, stateDataKey)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, fsm.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get state data: %w", err)
	}

	// unmarshal state data
	var stateData map[string]interface{}
	if err := json.Unmarshal(stateDataBytes, &stateData); err != nil {
		return nil, fmt.Errorf("unmarshal state data: %w", err)
	}

	// get data
	res, ok := stateData[key]
	if !ok {
		return nil, fsm.ErrNotFound
	}

	return res, nil
}

func (s *Storage) Close() error {
	return s.rds.Close()
}

func (s *Storage) generateKey(chat, user int64, keyType keyType) string {
	return fmt.Sprintf("%s:%d:%d:%s", s.prefix, chat, user, keyType)
}
