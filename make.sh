#!/bin/bash
set -e

function print_help {
	printf "Available Commands:\n";
	awk -v sq="'" '/^function run_([a-zA-Z0-9-]*)\s*/ {print "-e " sq NR "p" sq " -e " sq NR-1 "p" sq }' "$0" \
		| while read line; do eval "sed -n $line $0"; done \
		| paste -d"|" - - \
		| sed -e 's/^/  /' -e 's/function run_//' -e 's/#//' -e 's/{/	/' \
		| awk -F '|' '{ print "  " $2 "\t" $1}' \
		| expand -t 30
}

function run_test { # test the codebase
	go test -covermode=count -coverprofile=/tmp/cover \
		./ \
		./coding \
		./accept \
		&& go tool cover -html=/tmp/cover 

	go test -race -test.run=TestServerHandling
	go test -test.bench=.* -benchmem
}

case $1 in
	"test") run_test ;;
	*) print_help ;;
esac