#!/bin/bash

# Quip MCP Integration Test Runner
# This script provides a convenient way to run various types of tests

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}=== $1 ===${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

check_token() {
    if [ -z "$QUIP_API_TOKEN" ]; then
        print_error "QUIP_API_TOKEN environment variable is required for integration tests"
        echo "Set it with: export QUIP_API_TOKEN=your-token-here"
        echo "Get your token from: https://quip.com/dev/token"
        exit 1
    fi
    print_success "QUIP_API_TOKEN is set"
}

run_unit_tests() {
    print_header "Running Unit Tests"
    go test -v ./...
    print_success "Unit tests completed"
}

run_integration_tests() {
    print_header "Running Integration Tests"
    check_token
    go test -v -tags=integration ./pkg/quip -run TestIntegration
    print_success "Integration tests completed"
}

run_specific_test() {
    local test_name=$1
    if [ -z "$test_name" ]; then
        print_error "Test name is required"
        echo "Usage: $0 single <TestName>"
        echo "Example: $0 single GetRecentThreads"
        exit 1
    fi
    
    print_header "Running Specific Integration Test: $test_name"
    check_token
    go test -v -tags=integration ./pkg/quip -run TestIntegration_$test_name
    print_success "Test $test_name completed"
}

run_benchmarks() {
    print_header "Running Integration Benchmarks"
    check_token
    go test -v -tags=integration -bench=BenchmarkIntegration ./pkg/quip
    print_success "Benchmarks completed"
}

run_coverage() {
    print_header "Running Tests with Coverage"
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    print_success "Coverage report generated: coverage.html"
}

debug_recent_threads() {
    print_header "Debugging Recent Threads Issue"
    check_token
    print_warning "This will show debug output from the GetRecentThreads API call"
    go test -v -tags=integration ./pkg/quip -run TestIntegration_GetRecentThreads
}

show_help() {
    echo "Quip MCP Test Runner"
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  unit              Run unit tests only"
    echo "  integration       Run integration tests (requires QUIP_API_TOKEN)"
    echo "  single <test>     Run specific integration test"
    echo "  benchmarks        Run integration benchmarks"
    echo "  coverage          Run tests with coverage report"
    echo "  all               Run all tests (unit + integration)"
    echo "  debug-threads     Debug recent threads API issue"
    echo "  help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 unit"
    echo "  $0 integration"
    echo "  $0 single GetRecentThreads"
    echo "  $0 benchmarks"
    echo "  $0 debug-threads"
    echo ""
    echo "Environment:"
    echo "  QUIP_API_TOKEN    Required for integration tests"
    echo "                    Get yours from: https://quip.com/dev/token"
}

# Main script logic
case "${1:-help}" in
    "unit")
        run_unit_tests
        ;;
    "integration")
        run_integration_tests
        ;;
    "single")
        run_specific_test "$2"
        ;;
    "benchmarks")
        run_benchmarks
        ;;
    "coverage")
        run_coverage
        ;;
    "all")
        run_unit_tests
        echo ""
        if [ -n "$QUIP_API_TOKEN" ]; then
            run_integration_tests
        else
            print_warning "Skipping integration tests (QUIP_API_TOKEN not set)"
        fi
        ;;
    "debug-threads")
        debug_recent_threads
        ;;
    "help")
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac
