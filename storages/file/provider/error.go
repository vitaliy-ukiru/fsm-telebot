package provider

import "github.com/vitaliy-ukiru/fsm-telebot/storages/file"

func newError(provider string, op string, err error) error {
	if err == nil {
		return nil
	}
	return &file.ProviderError{ProviderType: provider, Operation: op, Err: err}
}
