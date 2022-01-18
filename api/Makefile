
.PHONY: fmt
fmt: ## Run go fmt for the whole project
	test -z $$(for d in $$(go list -f {{.Dir}} ./...); do gofmt -e -l -w $$d/*.go; done)
