plugins_dir := ~/.op/plugins/local

.PHONY: new-plugin registry %/example-secrets %/validate %/build test

beta-notice:
	@echo "# BETA NOTICE: The plugin ecosystem is in beta and is subject to change."
	@echo "# You may have to change or recompile your local builds every now and then to keep it"
	@echo "# compatible with the 1Password CLI updates."
	@echo

new-plugin: beta-notice
	go run cmd/contrib/main.go $@

registry:
	@rm -f plugins/plugins.go
	@go run cmd/contrib/main.go $@

%/example-secrets: registry
	go run cmd/contrib/main.go $@

%/validate: registry beta-notice
	go run cmd/contrib/main.go $@

validate: registry
	go run cmd/contrib/main.go $@

$(plugins_dir):
	mkdir -p $(plugins_dir)
	chmod 700 $(plugins_dir)

%/build: $(plugins_dir) registry beta-notice
	$(eval plugin := $(firstword $(subst /, ,$@)))
	@go run cmd/contrib/main.go $(plugin)/exists
	go build -o $(plugins_dir)/$(plugin) -ldflags="-X 'main.PluginName=$(plugin)'" ./cmd/contrib/build/

test:
	go test ./...
