# check to see if we can use ash, in alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)


run:
	# The below pipes the first program's output towards stdOut to second programs StdIn
	go run apis/services/sales/main.go | go run apis/tooling/logfmt/main.go

tidy:
	go mod tidy
	# This is putting all third part code packages in the vendor folder.
	go mod vendor
