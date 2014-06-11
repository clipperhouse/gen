@echo off
go get
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
