#!/bin/bash

echo "Building Prime XLSX Monthly Summary Reporting Go app for Linux (amd64)..."

GOOS=linux GOARCH=amd64 go build -o ./_bin/prime-monthly-summary-reporting ./cmd

if [ $? -eq 0 ]; then
    echo "Build successful! Output: prime-monthly-summary-reporting"
else
    echo "Build failed!"
fi
