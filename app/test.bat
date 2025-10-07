@echo off
REM Batch file for running tests on Windows

if "%1"=="" goto help
if "%1"=="help" goto help
if "%1"=="test" goto test
if "%1"=="test-verbose" goto test-verbose
if "%1"=="test-coverage" goto test-coverage
if "%1"=="test-services" goto test-services
if "%1"=="build" goto build
if "%1"=="run" goto run
if "%1"=="clean" goto clean
goto help

:help
echo Available commands:
echo   test.bat test          - Run all tests
echo   test.bat test-verbose  - Run all tests with verbose output
echo   test.bat test-coverage - Run tests with coverage report
echo   test.bat test-services - Run service layer tests only
echo   test.bat build         - Build the application
echo   test.bat run           - Run the application
echo   test.bat clean         - Clean build artifacts
goto end

:test
echo 🧪 Running all tests...
go test ./tests ./tests/services
goto end

:test-verbose
echo 🧪 Running all tests (verbose)...
go test -v ./tests ./tests/services
goto end

:test-coverage
echo 📊 Running tests with coverage...
go test -coverprofile=coverage.out ./tests ./tests/services
go tool cover -html=coverage.out -o coverage.html
echo 📄 Coverage report generated: coverage.html
goto end

:test-services
echo 🔧 Running service tests...
go test -v ./tests/services
goto end

:build
echo 🔨 Building application...
if not exist bin mkdir bin
go build -o bin/app.exe main.go
goto end

:run
echo 🚀 Starting application...
go run main.go
goto end

:clean
echo 🧹 Cleaning up...
if exist coverage.out del coverage.out
if exist coverage.html del coverage.html
if exist bin rmdir /s /q bin
goto end

:end
