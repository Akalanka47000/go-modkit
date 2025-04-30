format:
	for mod in $(find . -name go.mod); do \
		dir=$(dirname "$mod") \
		cd "$dir" && go fmt ./... \
	done

test:
	go test -v ./...

lint:
	golangci-lint run ./...