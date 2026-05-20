package limiter

import (
	"time"

	"backend-assignment/internal/storage"
)

const MaxRequests = 5
const WindowDuration = time.Minute

func AllowRequest(userID string) bool {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	userData, exists := storage.RequestStore[userID]

	if !exists {
		userData = &storage.UserRequestData{}
		storage.RequestStore[userID] = userData
	}

	now := time.Now()

	var validTimestamps []time.Time

	for _, ts := range userData.AcceptedTimestamps {
		if now.Sub(ts) < WindowDuration {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	userData.AcceptedTimestamps = validTimestamps

	if len(userData.AcceptedTimestamps) >= MaxRequests {
		userData.RejectedCount++
		return false
	}

	userData.AcceptedTimestamps = append(userData.AcceptedTimestamps, now)
	userData.AcceptedCount++

	return true
}
