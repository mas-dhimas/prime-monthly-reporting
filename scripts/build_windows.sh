#!/bin/bash

echo "Building Prime XLSX Monthly Summary Reporting Go app for Windows (amd64)..."

GOOS=windows GOARCH=amd64 go build -o ./_bin/prime-monthly-summary-reporting.exe ./cmd

if [ $? -eq 0 ]; then
    echo "Build successful! Output: prime-xlsx-monthly-summary-reporting.exe"
else
    echo "Build failed!"
fi