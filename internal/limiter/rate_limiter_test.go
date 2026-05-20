package limiter

import (
	"testing"

	"backend-assignment/internal/storage"
)

func TestAllowRequest_LimitEnforcement(t *testing.T) {

	userID := "test-user"

	// Reset test state
	storage.RequestStore = make(map[string]*storage.UserRequestData)

	// First 5 requests should pass
	for i := 0; i < 5; i++ {

		allowed := AllowRequest(userID)

		if !allowed {
			t.Errorf("expected request %d to be allowed", i+1)
		}
	}

	// 6th request should fail
	allowed := AllowRequest(userID)

	if allowed {
		t.Errorf("expected 6th request to be rejected")
	}
}
