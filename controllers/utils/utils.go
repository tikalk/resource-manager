package utils

import (
	"fmt"
	"time"
)

// IsObjExpired check if object has expired
func IsObjExpired(creation time.Time, expiration string) (err error, secondsUntilExp int) {
	now := time.Now()
	objCreatedAt := creation                           // creation time
	expireAfter, err := time.ParseDuration(expiration) // will be expired after (object life period)
	objExpiredAt := objCreatedAt.Add(expireAfter)      // calculate the exact time it will expire
	// calculate how many seconds has left until expiration date
	secondsUntilExp = int(objExpiredAt.Sub(now).Seconds())
	return err, secondsUntilExp
}

func IsIntervalOccurred(now time.Time, timeframe string) (err error, secondsUntilTimeframe int) {
	// parse timeframe to time object
	timeframeTime, err := time.Parse("15:04", timeframe)
	if err != nil {
		fmt.Println("Could not parse timeframe:", err)
		return err, 0
	}

	// current hour in the same time format "15:04"
	nowTime, err := time.Parse("15:04", now.Format("15:04"))
	if err != nil {
		fmt.Println("Could not parse time:", err)
		return err, 0
	}

	secondsUntilTimeframe = int(timeframeTime.Sub(nowTime).Seconds())
	//if secondsUntilTimeframe <= 0 && secondsUntilTimeframe > -60 {
	//	fmt.Println(fmt.Sprintf("timeframe triggerd: '%s'", timeframeTime))
	//	return nil, true
	//}
	return nil, secondsUntilTimeframe
}
