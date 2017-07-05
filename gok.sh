#!/bin/sh
#
# Verifies that go code passes go fmt, go vet, golint, and go test.
#

o=$(mktemp tmp.XXXXXXXXXX)

fail() {
	echo Failed
	cat $o
	rm $o
	exit 1
}

trap fail INT TERM

#echo Generating
#go generate ./parser/camxes || fail
#go generate ./parser/ilmentufa || fail
#go generate ./parser/maftufa || fail

echo Formatting
gofmt -l $(find . -name '*.go') > $o 2>&1 
test $(wc -l $o | awk '{ print $1 }') = "0" || fail

echo Vetting
go vet ./... > $o 2>&1 || fail

echo Testing
go test -test.timeout=10s ./... > $o 2>&1 || fail

echo Linting
golint ./... \
	| egrep -v "parser/.*\.go.*don't use underscores"\
	| egrep -v "parser/.*\.go.*ALL_CAPS"\
	| egrep -v "parser/.*/all\.go.*a blank import"\
	> $o 2>&1
# Silly: diff the grepped golint output with empty.
# If it's non-empty, error, otherwise succeed.
e=$(tempfile)
touch $e
diff $o $e > /dev/null || { rm $e; fail; }

rm $o $e
