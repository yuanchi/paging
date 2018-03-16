package util

import (
	"time"
	"log"	
)

var taipei *time.Location
const defaultFormat = "20060102-150405.000000-"

func init() {
	timezone, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Fatalln(err)	
	}
	taipei = timezone
}

func deinit() {
	taipei = nil
}

func CurrentTimef() string {
	now := time.Now().In(taipei)
	f := now.Format(defaultFormat)
	return f
}
