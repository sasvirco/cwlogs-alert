module process

go 1.14

replace github.com/cwlogsalert/model => ../model

require (
	github.com/aws/aws-sdk-go v1.34.25 // indirect
	github.com/cwlogsalert/model v0.0.0-00010101000000-000000000000
	github.com/hashicorp/go-memdb v1.2.1
	github.com/rs/zerolog v1.20.0
)
