module src/main.go

go 1.14

require (
	github.com/aws/aws-sdk-go v1.34.27
	github.com/cwlogsalert/config v0.0.0-00010101000000-000000000000
	github.com/cwlogsalert/cwlog v0.0.0-00010101000000-000000000000
	github.com/cwlogsalert/db v0.0.0-00010101000000-000000000000
	github.com/cwlogsalert/model v0.0.0-00010101000000-000000000000
	github.com/cwlogsalert/notify v0.0.0-00010101000000-000000000000
	github.com/cwlogsalert/process v0.0.0-00010101000000-000000000000
	github.com/pelletier/go-toml v1.8.1
	github.com/rs/zerolog v1.20.0
)

replace (
	github.com/cwlogsalert/config => ./src/config
	github.com/cwlogsalert/cwlog => ./src/cwlog
	github.com/cwlogsalert/db => ./src/db
	github.com/cwlogsalert/model => ./src/model
	github.com/cwlogsalert/notify => ./src/notify
	github.com/cwlogsalert/process => ./src/process
)
