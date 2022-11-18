plugins_dir := ~/.op/plugins/local

new-plugin:
	go run cmd/contrib/main.go plugin

%/example-secrets:
	go run cmd/contrib/main.go $@

%/validate:
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