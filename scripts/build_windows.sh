#!/bin/bash

echo "Building Prime XLSX Monthly Summary Reporting Go app for Windows (amd64)..."

GOOS=windows GOARCH=amd64 go build -o prime-xslx-monthly-summary-reporting.exe

if [ $? -eq 0 ]; then
    echo "Build successful! Output: prime-xlsx-monthly-summary-reporting.exe"
else
    echo "Build failed!"
fi