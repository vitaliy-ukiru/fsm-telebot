// Package storages contains implementations of state storage.
//
// # Memory storage
//
// This is concurrency safe in-memory storage.
// It'll be reset on shutdown.
//
// # Redis storage
//
// This storage based in redis.
//
// Distribute on separate module as
// [github.com/nacknime-official/fsm-telebot-redis-storage/v2]
//
// All issues/pull request send to this [repository].
// Please note that this repository separates major
// versions by directory.
//
// [repository]: https://github.com/nacknime-official/fsm-telebot-redis-storage
package storage
