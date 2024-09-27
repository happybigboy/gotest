// Package states manages the user states for the Telegram bot
package states

import (
	"log"
	"sync"
)

type UserState struct {
	mu     sync.Mutex
	states map[int64]string
}

// NewUserState initializes the UserState struct
func NewUserState() *UserState {
	log.Printf("New state")
	return &UserState{
		states: make(map[int64]string),
	}
}

// SetState sets the state for a specific chat ID in a thread-safe way
func (us *UserState) SetState(chatID int64, state string) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.states[chatID] = state
	log.Printf("Set state for chatID %d: %s", chatID, state)
}

// GetState retrieves the state for a specific chat ID in a thread-safe way
func (us *UserState) GetState(chatID int64) string {
	us.mu.Lock()
	defer us.mu.Unlock()
	log.Printf("Get state for chatID %d: %s", chatID, us.states[chatID])
	return us.states[chatID]

}

// ResetState deletes the state for a specific chat ID in a thread-safe way
func (us *UserState) ResetState(chatID int64) {
	log.Printf("DELETE")
	us.mu.Lock()
	defer us.mu.Unlock()
	delete(us.states, chatID)
}
