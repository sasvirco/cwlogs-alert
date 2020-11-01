module cwlog

go 1.14

require (
	github.com/aws/aws-sdk-go v1.34.22
	github.com/cwlogsalert/config v0.0.0-00010101000000-000000000000
	github.com/cwlogsalert/model v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.19.0
)

replace (
	github.com/cwlogsalert/config => ../config
	github.com/cwlogsalert/model => ../model
)
