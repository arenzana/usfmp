# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Updated GitHub Actions dependencies:
  - golangci/golangci-lint-action from v6 to v8
  - anchore/scan-action from v4 to v6  
  - actions/download-artifact from v4 to v5
- Updated Go dependencies:
  - github.com/spf13/pflag from v1.0.6 to v1.0.7

### Planned
- PDF output formatter implementation
- Enhanced error reporting with line-by-line context
- Support for additional USFM markers (poetry, lists, etc.)
- Configuration file support for CLI
- Parallel processing for directory parsing
- Plugin system for custom formatters

## [0.0.2] - 2024-01-XX

### Added
- Initial release of USFM Parser
- Complete USFM 3.1 parser implementation
- Support for major USFM markers:
  - Document identification (`\id`, `\h`, `\toc1-3`, `\mt1`)
  - Chapter and verse markers (`\c`, `\v`)
  - Section headers (`\s1`, `\s2`, `\s3`) with multiple levels
  - Cross-references (`\r`)
  - Footnotes (`\f`, `\fr`, `\ft`, `\f*`) with extraction and cleaning
- Multiple output formatters:
  - JSON formatter with single document and array support
  - Plain text formatter with human-readable output
  - TSV formatter for data analysis
  - PDF formatter placeholder
- Command-line interface with Cobra framework:
  - Single file and directory processing
  - Multiple output formats (`-f json|txt|tsv|pdf`)
  - Verbosity control (`--verbose`, `--quiet`)
  - Strict parsing mode (`--strict`)
  - Output file specification (`--output`)
- Go library with comprehensive API:
  - `Parser` type with configurable options
  - `Document` hierarchy (Document → Chapter → Section → Verse)
  - `ParseOptions` for controlling parser behavior
  - Rich type definitions with JSON serialization support
- Parser features:
  - Strict vs lenient parsing modes
  - Footnote extraction and text cleaning
  - Cross-reference processing
  - Section hierarchy support
  - Comprehensive error handling
- Development infrastructure:
  - Comprehensive Makefile with multiple targets
  - Full unit test suite with >95% coverage
  - Integration tests with real biblical text samples
  - Code formatting and linting support
  - Multi-platform build support
- Documentation:
  - Complete GoDoc documentation for all public APIs
  - Comprehensive README with usage examples
  - CLAUDE.md for AI assistant guidance
  - Sample data integration (Berean Standard Bible)

### Technical Details
- Written in Go with standard library focus
- Uses regular expressions for efficient marker parsing
- Supports concurrent parsing preparation
- Memory-efficient streaming parser design
- Comprehensive error handling with context
- Clean separation of concerns (parser, formatters, CLI)

### Sample Data
- Included complete Berean Standard Bible in USFM format
- 66 biblical books for comprehensive testing
- Real-world footnotes and cross-references
- Various section structures and complexity levels

### Dependencies
- [cobra](https://github.com/spf13/cobra) v1.9.1 - CLI framework
- Go 1.24+ standard library

### Build System
- Make-based build system with comprehensive targets
- Cross-platform compilation support (Linux, macOS, Windows)
- Automated testing and coverage reporting
- Code quality tools integration
- Development workflow automation

### File Structure
```
usfmp/
├── cmd/usfmp/           # CLI application
├── pkg/usfm/            # Public library API
├── internal/formatter/  # Output formatters  
├── bsb_usfm/           # Sample biblical texts
├── build/              # Build outputs
└── docs/               # Additional documentation
```

---

## Release History

### Version 0.0.2 - Initial Release
- **Release Date**: 2024-01-XX
- **Go Version**: 1.24+
- **Breaking Changes**: None (initial release)
- **Migration Guide**: N/A (initial release)

## Development Workflow

### Versioning Strategy
This project follows [Semantic Versioning](https://semver.org/):
- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality in a backwards compatible manner  
- **PATCH** version for backwards compatible bug fixes

### Release Process
1. Update CHANGELOG.md with new version details
2. Update version in Makefile and relevant files
3. Create and test release candidate
4. Run full test suite across all supported platforms
5. Create Git tag with version number
6. Build and publish release artifacts
7. Update documentation and examples

### Breaking Changes Policy
Breaking changes will be:
- Clearly documented in CHANGELOG
- Include migration guide when applicable
- Follow deprecation warnings when possible
- Announced in advance for major changes

### Support Policy
- **Latest Major Version**: Full support with new features and bug fixes
- **Previous Major Version**: Security fixes and critical bug fixes only
- **Older Versions**: No active support (community patches accepted)

---

## Contributors

### Core Team
- [@arenzana](https://github.com/arenzana) - Project Creator & Lead Developer

### Special Thanks
- USFM Working Group for the specification
- Berean Bible for providing sample texts
- Go community for excellent tooling and libraries

---

*This changelog is maintained by the project maintainers. For automated changelog generation, see the [GitHub Releases](https://github.com/arenzana/usfmp/releases) page.*