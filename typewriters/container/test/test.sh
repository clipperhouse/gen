go get -u
go run setup.go
touch coverage.out
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
