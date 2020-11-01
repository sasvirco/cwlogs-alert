package model

import (
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/cwlogsalert/config"
)

//ResultItem - channel message structure
type ResultItem struct {
	Result *cloudwatchlogs.GetQueryResultsOutput
	Rule   config.Rule
}

//NotificationItem - channel message structure
type NotificationItem struct {
	Rule   config.Rule
	Result *cloudwatchlogs.GetQueryResultsOutput
	State  string
}

//Alert -
type Alert struct {
	RuleName string
	Alerting bool
}

//MakeTimestamp - return now()
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//ConvertDurationToMilliseconds - converts string representation of duration to milliseconds
func ConvertDurationToMilliseconds(duration string) int64 {
	u, _ := time.ParseDuration(duration)
	return u.Milliseconds()
}
