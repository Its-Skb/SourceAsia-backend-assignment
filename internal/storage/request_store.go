package storage

import (
	"sync"
	"time"
)

type UserRequestData struct {
	AcceptedTimestamps []time.Time
	AcceptedCount      int
	RejectedCount      int
}

var (
	RequestStore = make(map[string]*UserRequestData)
	Mu           sync.RWMutex
)
