#!/usr/bin/env bash

# Check if the correct number of arguments are provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <expected_coverage>"
    exit 1
fi

EXPECTED_COVERAGE="$1"

rm -rf coverage
mkdir coverage

go test ./... -v -coverprofile coverage/all

go tool cover -html=coverage/all -o coverage/index.html
go tool cover -func=coverage/all > coverage/output

ACTUAL_CORVERAGE=$(grep -E 'total:\s+\(statements\)\s+(\d{1,2}\.\d{1,2})%' coverage/output | sed -n 's/.*\([0-9]\{1,2\}\.[0-9]\{1,2\}\)%.*/\1/p')

result=$(echo "$ACTUAL_CORVERAGE >= $EXPECTED_COVERAGE" | bc)

if [ "$result" -eq 1 ]; then
  exit 0
else
  echo "Actual coverage ($ACTUAL_CORVERAGE) is less than expected coverage ($EXPECTED_COVERAGE)."
  exit 1
fi