go test setup_test.go >/dev/null
touch coverage.out
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
