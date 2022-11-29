plugins_dir := ~/.op/plugins/local

new-plugin:
	go run cmd/contrib/main.go $@

registry:
	@rm -f plugins/plugins.go
	go run cmd/contrib/main.go $@

%/example-secrets: registry
	go run cmd/contrib/main.go $@

%/validate: registry
	go run cmd/contrib/main.go $@

validate: registry
	go run cmd/contrib/main.go $@

$(plugins_dir):
	mkdir -p $(plugins_dir)
	chmod 700 $(plugins_dir)

%/build: $(plugins_dir)
	$(eval plugin := $(firstword $(subst /, ,$@)))
	@go run cmd/contrib/main.go $(plugin)/exists
	go build -o $(plugins_dir)/$(plugin) -ldflags="-X 'main.PluginName=$(plugin)'" ./cmd/contrib/build/

test:
	go test ./...
