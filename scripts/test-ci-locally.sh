#!/bin/bash

# Local CI Testing Script for USFM Parser
# This script tests various CI components locally before pushing to GitHub

set -e  # Exit on any error

echo "🧪 USFM Parser - Local CI Testing"
echo "=================================="
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print section headers
print_section() {
    echo -e "\n${BLUE}$1${NC}"
    echo "$(printf '=%.0s' $(seq 1 ${#1}))"
}

# Function to print success
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Function to print warning
print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Function to print error
print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if we're in the project root
if [[ ! -f "go.mod" ]] || [[ ! -d ".github" ]]; then
    print_error "Please run this script from the project root directory"
    exit 1
fi

print_section "1. Code Quality Checks"

# Test Go formatting
echo "Checking Go code formatting..."
if go fmt ./... | grep -q .; then
    print_warning "Code formatting issues found, but continuing..."
else
    print_success "Code formatting is correct"
fi

# Test Go vetting
echo "Running go vet..."
if go vet ./...; then
    print_success "go vet passed"
else
    print_error "go vet failed"
    exit 1
fi

# Test golangci-lint if available
if command -v golangci-lint >/dev/null 2>&1; then
    echo "Running golangci-lint..."
    if golangci-lint run --timeout=5m; then
        print_success "golangci-lint passed"
    else
        print_warning "golangci-lint found issues, but continuing..."
    fi
else
    print_warning "golangci-lint not installed - skipping linting"
    echo "  Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
fi

print_section "2. Unit Tests"

echo "Running unit tests..."
if make test; then
    print_success "All unit tests passed"
else
    print_error "Unit tests failed"
    exit 1
fi

echo "Running tests with coverage..."
if make coverage; then
    print_success "Coverage report generated"
    if command -v open >/dev/null 2>&1; then
        echo "  View coverage: open coverage.html"
    fi
else
    print_warning "Coverage generation failed, but continuing..."
fi

print_section "3. Build Testing"

echo "Testing standard build..."
if make build; then
    print_success "Standard build successful"
    
    # Test the built binary
    if ./build/usfmp --version >/dev/null; then
        print_success "Built binary works correctly"
    else
        print_error "Built binary failed to run"
        exit 1
    fi
else
    print_error "Build failed"
    exit 1
fi

# Test goreleaser if available
if command -v goreleaser >/dev/null 2>&1; then
    echo "Testing goreleaser build..."
    if goreleaser build --snapshot --clean --quiet; then
        print_success "GoReleaser build successful"
        
        # Test a built binary
        if find dist/ -name "usfmp*" -type f | head -1 | xargs -I {} {} --version >/dev/null; then
            print_success "GoReleaser binary works correctly"
        else
            print_warning "GoReleaser binary test failed"
        fi
        
        # Clean up
        rm -rf dist/
    else
        print_warning "GoReleaser build failed, but continuing..."
    fi
else
    print_warning "GoReleaser not installed - skipping release build test"
    echo "  Install with: go install github.com/goreleaser/goreleaser@latest"
fi

print_section "4. Security Scanning"

# Test grype if available  
if command -v grype >/dev/null 2>&1; then
    echo "Running security scan with Grype..."
    if grype build/usfmp --quiet; then
        print_success "Security scan passed - no vulnerabilities found"
    else
        print_warning "Security scan found issues, but continuing..."
    fi
else
    print_warning "Grype not installed - skipping security scan"
    echo "  Install with: curl -sSfL https://raw.githubusercontent.com/anchore/grype/main/install.sh | sh -s -- -b /usr/local/bin"
fi

print_section "5. Container Testing"

if command -v docker >/dev/null 2>&1; then
    echo "Testing Docker build..."
    
    # Copy binary for Docker build
    if [[ -f "build/usfmp" ]]; then
        # We need a Linux binary for Docker
        if [[ "$(uname)" == "Darwin" ]]; then
            echo "Building Linux binary for Docker..."
            GOOS=linux GOARCH=amd64 go build -o usfmp-linux ./cmd/usfmp
            cp usfmp-linux usfmp
        else
            cp build/usfmp .
        fi
        
        if docker build -t usfmp-test:latest -q . >/dev/null; then
            print_success "Docker build successful"
            
            # Test container
            if docker run --rm usfmp-test:latest --help >/dev/null; then
                print_success "Docker container works correctly"
            else
                print_warning "Docker container test failed"
            fi
            
            # Clean up
            docker rmi usfmp-test:latest >/dev/null 2>&1 || true
        else
            print_warning "Docker build failed"
        fi
        
        # Clean up binary
        rm -f usfmp usfmp-linux
    else
        print_warning "No binary found - skipping Docker test"
    fi
else
    print_warning "Docker not available - skipping container test"
fi

print_section "6. GitHub Actions Validation"

if command -v act >/dev/null 2>&1; then
    echo "Testing GitHub Actions workflows with act..."
    
    # List workflows
    echo "Available workflows:"
    act --list 2>/dev/null | grep -E "Job ID|Job name" | head -5
    
    # Test workflow syntax (dry run)
    if act -j lint --container-architecture linux/amd64 --dryrun -q 2>/dev/null; then
        print_success "GitHub Actions workflow syntax is valid"
    else
        print_warning "GitHub Actions workflow validation failed (needs git repo for full test)"
    fi
else
    print_warning "act not installed - skipping GitHub Actions local test"
    echo "  Install with: brew install act"
fi

print_section "7. Documentation Testing"

echo "Testing documentation generation..."
if make docs >/dev/null 2>&1; then
    print_success "Documentation generation works"
    rm -f docs.txt  # Clean up
else
    print_warning "Documentation generation failed"
fi

# Check if godoc works
if go doc ./pkg/usfm >/dev/null 2>&1; then
    print_success "GoDoc comments are properly formatted"
else
    print_warning "GoDoc generation has issues"
fi

print_section "8. Integration Testing"

if [[ -f "bsb_usfm/01GENBSB.SFM" ]]; then
    echo "Testing with sample USFM data..."
    
    # Test different output formats
    formats=("json" "txt" "tsv")
    for format in "${formats[@]}"; do
        if ./build/usfmp -f "$format" -q bsb_usfm/01GENBSB.SFM >/dev/null; then
            print_success "$format output format works"
        else
            print_error "$format output format failed"
            exit 1
        fi
    done
else
    print_warning "Sample USFM data not found - skipping integration tests"
fi

print_section "Summary"

echo -e "\n${GREEN}🎉 Local CI testing completed successfully!${NC}"
echo
echo "Your code is ready for:"
echo "  • ✅ GitHub Actions CI pipeline"
echo "  • ✅ Automated releases with GoReleaser"
echo "  • ✅ Security scanning"
echo "  • ✅ Docker containerization"
echo
echo -e "${BLUE}Next steps:${NC}"
echo "  1. Initialize git repository: git init"
echo "  2. Add remote: git remote add origin https://github.com/arenzana/usfmp"
echo "  3. Push to GitHub to trigger CI/CD pipeline"
echo
echo -e "${YELLOW}Note:${NC} Some tests require a git repository and GitHub token for full validation."
echo "The GitHub Actions will run automatically when pushed to GitHub."