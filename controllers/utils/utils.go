package utils

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// IsObjExpired check if object has expired
func IsObjExpired(creation v1.Time, expiration string) (bool, int) {

	now := time.Now()
	objCreatedAt := creation                         // creation time
	expireAfter, _ := time.ParseDuration(expiration) // will be expired after (object life period)
	objExpiredAt := objCreatedAt.Add(expireAfter)    // calculate the exact time it will expire
	// calculate how many seconds has left until expiration date
	secondsUntilExp := objExpiredAt.Sub(now).Seconds()
	return secondsUntilExp <= 0, int(secondsUntilExp)
}
