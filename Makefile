# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

# Add the ability to override some variables
# Use with care
-include override.mk

# Main targets
include main.mk

# Add custom targets here
-include custom.mk

.PHONY: generate
generate: ## Generate test code
	go run sagikazarmark.dev/mga generate kit endpoint ./...
	go run sagikazarmark.dev/mga generate event handler ./...
	go run sagikazarmark.dev/mga generate event handler --output subpkg:suffix=gen ./...
	go run sagikazarmark.dev/mga generate event dispatcher ./...
	go run sagikazarmark.dev/mga generate event dispatcher --output subpkg:suffix=gen ./...
