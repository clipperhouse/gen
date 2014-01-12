#!/bin/bash

# Delete generated file if they exists.
if [ -f movie_gen.go ]; then
	rm movie_gen.go
fi

# Build `gen` binary according to current source,
# regenerate the `Movie` struct and remove
# the binary.
go build -o test_gen ..
./test_gen -f *models.Movie # Need to force, otherwise movie_test.go doesn't compile and blocks `gen`
rm test_gen

# Run the test with the updated `Movies` struct.
go test .
