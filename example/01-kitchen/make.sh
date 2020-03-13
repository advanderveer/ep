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

function run_vendor { # get the dependencies for this example
	rm -fr vendor
	git clone https://github.com/gorilla/mux.git vendor/github.com/gorilla/mux 
	git clone https://github.com/alexedwards/scs.git vendor/github.com/alexedwards/scs/v2 
	git clone https://github.com/go-playground/form.git vendor/github.com/go-playground/form/v4
	mkdir -p vendor/github.com/advanderveer
	ln -s $(cd ../..; pwd) vendor/github.com/advanderveer/ep
}

function run_run { # run the example
	go run .
}

case $1 in
	"vendor") run_vendor ;;
	"run") run_run ;;
	*) print_help ;;
esac