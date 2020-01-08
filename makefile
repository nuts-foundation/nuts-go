update-nuts-deps:
	cat go.mod | awk '/nuts-foundation.* / {print $$1}' | xargs go get

test:
	go test ./...

