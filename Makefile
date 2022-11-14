new-plugin:
	go run cmd/contrib/main.go plugin

%/example-secrets:
	go run cmd/contrib/main.go $@

%/validate:
	go run cmd/contrib/main.go $@

test:
	go test ./...