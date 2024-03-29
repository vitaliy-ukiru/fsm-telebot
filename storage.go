package fsm

import "errors"

// ErrNotFound returns if data not found.
var ErrNotFound = errors.New("fsm/storage: not found")

// Storage is object what uses for save information for FSM.
// It can be client for DBMS, file or just in memory storage.
//
// In package storages you can find some implementation.
//
// You can contribute your implementations to pull requests
// or create your repository.
//
// Not recommended works with storage from handlers.
type Storage interface {
	// GetState returns State for target. Default state
	// is empty string (DefaultState).
	GetState(chatId, userId int64) (State, error)

	// SetState sets state for target.
	SetState(chatId, userId int64, state State) error

	// ResetState deletes state for target. If `withData` is true
	// deletes user data from storage.
	ResetState(chatId, userId int64, withData bool) error

	// UpdateData sets, updates or deletes data for target. When
	//data argument is nil it must deletes this item.
	UpdateData(chatId, userId int64, key string, data any) error

	// GetData gets data for target and saves it into `to` argument.
	// Destination argument must be a valid pointer.
	GetData(chatId, userId int64, key string, to any) error

	// Close closes storage. Needs for correct work
	// with storage connection.
	Close() error
}
