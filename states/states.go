// Package states manages the user states for the Telegram bot
package states

import (
	"sync"
)

type UserState struct {
	mu     sync.Mutex
	states map[int64]string
}

// NewUserState initializes the UserState struct
func NewUserState() *UserState {
	return &UserState{
		states: make(map[int64]string),
	}
}

// SetState sets the state for a specific chat ID in a thread-safe way
func (us *UserState) SetState(chatID int64, state string) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.states[chatID] = state
}

// GetState retrieves the state for a specific chat ID in a thread-safe way
func (us *UserState) GetState(chatID int64) string {
	us.mu.Lock()
	defer us.mu.Unlock()
	return us.states[chatID]
}

// ResetState deletes the state for a specific chat ID in a thread-safe way
func (us *UserState) ResetState(chatID int64) {
	us.mu.Lock()
	defer us.mu.Unlock()
	delete(us.states, chatID)
}