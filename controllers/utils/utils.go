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

	// current hour in 15:04 format
	currentHour := time.Now().Format("15:04")

	// parse timeframe to real time object
	timeframeTime, err := time.Parse("15:04", timeframe)
	// current hour in the same time format
	nowTime, err := time.Parse("15:04", currentHour)
	if err != nil {
		fmt.Println("Could not parse time:", err)
		return err, false
	}
	fmt.Println(fmt.Sprintf("The time now is: %s", nowTime))
	fmt.Println(fmt.Sprintf("The timeframe is: %s", timeframeTime))
	secondsUntilInterval := timeframeTime.Sub(nowTime).Seconds()
	fmt.Println(fmt.Sprintf("secondsUntilInterval: %d", int(secondsUntilInterval)))
	if secondsUntilInterval <= 0 {
		fmt.Println(fmt.Sprintf("timeframe triggerd: '%s'", timeframeTime))
		return nil, true
	}
	return nil, false
}
