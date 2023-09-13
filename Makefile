.PHONY: test race cover cover-out

test:
	go test

race:
	go test -race

cover:
	go test -cover

cover-out:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
