package utils

import (
	"fmt"
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

func IsIntervalOccurred(timeframe string) (error, bool) {

	// parse timeframe to time object
	timeframeTime, err := time.Parse("15:04", timeframe)
	if err != nil {
		fmt.Println("Could not parse timeframe:", err)
		return err, false
	}

	// current hour in the same time format "15:04"
	nowTime, err := time.Parse("15:04", time.Now().Format("15:04"))
	if err != nil {
		fmt.Println("Could not parse time:", err)
		return err, false
	}

	secondsUntilTimeframe := timeframeTime.Sub(nowTime).Seconds()
	if secondsUntilTimeframe <= 0 && secondsUntilTimeframe > -60 {
		fmt.Println(fmt.Sprintf("timeframe triggerd: '%s'", timeframeTime))
		return nil, true
	}
	return nil, false
}
