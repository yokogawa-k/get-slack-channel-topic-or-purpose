# vi: set ft=make ts=2 sw=2 sts=0 noet:

SHELL := /bin/bash

.PHONY: default
default: help

# http://postd.cc/auto-documented-makefile/
.PHONY: help help-common
help: help-common

help-common:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m %-30s\033[0m %s\n", $$1, $$2}'

build: ## Build
	DOCKER_BUILDKIT=1 docker build --output . .

exec: ## 実行
	./get-slacklog

# https://github.com/Respect/samples/blob/master/Makefile#L295
# Re-usable target for yes no prompt. Usage: make .prompt-yesno message="Is it yes or no?"
# Will exit with error if not yes
.PHONY: .prompt-yesno
.prompt-yesno:
	@exec 9<&0 0</dev/tty; \
	echo "$(message) [y/N]:"; \
	read -r -t 60 -n 3 yn; \
	exec 0<&9 9<&-; \
	if [[ -z $$yn ]]; then \
		echo "Please input y(es) or n(o)."; \
		exit 1; \
	else \
		if [[ $$yn =~ ^[yY] ]]; then \
			echo "continue" >&2; \
		else \
			echo "abort." >&2; \
			exit 1; \
		fi; \
	fi
