package fsm

import "github.com/pkg/errors"

// ErrNotFound returns if data not found.
var ErrNotFound = errors.New("fsm/storage: not found")

// Storage is object what uses for save information for FSM.
//
// It can be client for DBMS, file or just in memory storage (like storage.MemoryStorage)
// You can contribute your implementations to pull requests.
//
// For correct works as address of users uses chat id and user id (group support)
// Not recommended works with storage from handlers.
type Storage interface {
	// GetState returns State for target. Default state is empty string
	GetState(chatId, userId int64) (State, error)
	// SetState sets states for target.
	SetState(chatId, userId int64, state State) error
	// ResetState deletes state for target. If `withData` is true deletes user data from storage.
	ResetState(chatId, userId int64, withData bool) error

	// UpdateData sets, updates or deletes data for target. Set `data` as nil if you want delete.
	UpdateData(chatId, userId int64, key string, data interface{}) error

	// GetData gets data for target and saves it to `to`.
	// `to` must be a pointer.
	// If error is not nil then data will be nil.
	GetData(chatId, userId int64, key string, to interface{}) error

	// Close closes storage. Needs for correct work with storage connection.
	Close() error
}
