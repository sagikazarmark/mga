package pointer_message

import (
	"context"
)

type User struct{}
type TrackingInfo struct{}

type Service interface {
	// RegisterUser from https://github.com/sagikazarmark/mga/issues/32#issue-566519866
	RegisterUser(ctx context.Context, user User, trackingInfo *TrackingInfo) (registeredUser *User, r1 string, r2 string, err error)
}
