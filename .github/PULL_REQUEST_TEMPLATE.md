# Pull Request

## Description
A clear and concise description of what this PR does.

## Type of Change
Please mark the relevant option(s):

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)  
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] 📚 Documentation update
- [ ] 🔧 Refactoring (no functional changes, no api changes)
- [ ] ⚡ Performance improvement
- [ ] 🧪 Test addition or improvement
- [ ] 🔨 Build/CI changes

## Related Issue
Fixes # (issue number)
Relates to # (issue number)

## Changes Made
- List the main changes made in this PR
- Be specific about what was modified
- Include any new dependencies added

## Testing
Please describe the tests that you ran to verify your changes:

- [ ] Unit tests pass (`make test`)
- [ ] Integration tests with sample data pass
- [ ] Linting passes (`make lint`)
- [ ] Build succeeds (`make build`)
- [ ] Manual testing performed

### Test Configuration
- **Go version:** [e.g., 1.24]
- **OS:** [e.g., Linux, macOS, Windows]
- **Architecture:** [e.g., amd64, arm64]

### Manual Test Cases
If applicable, describe manual test cases:

```bash
# Example test command
usfmp -f json test-file.sfm
```

**Expected:** [describe expected outcome]
**Actual:** [describe actual outcome]

## Sample Input/Output
If this change affects parsing or output formatting, please provide examples:

### Input USFM
```usfm
\id GEN - Test
\c 1
\v 1 Test verse.
```

### Before (if applicable)
```json
{
  "error": "parsing failed"
}
```

### After
```json
{
  "id": "GEN - Test",
  "chapters": [...]
}
```

## Documentation
- [ ] Code comments updated/added
- [ ] README.md updated (if needed)
- [ ] CHANGELOG.md updated (will be updated by maintainer)
- [ ] GoDoc comments added/updated for public APIs
- [ ] Examples updated/added

## Breaking Changes
If this is a breaking change, please describe:

1. What breaks
2. Why this change was necessary
3. Migration guide for users
4. Deprecation timeline (if applicable)

## Performance Impact
- [ ] No performance impact
- [ ] Performance improvement (please quantify)
- [ ] Performance regression (please justify)

If there's a performance impact, please describe:
- What was measured
- Before/after metrics
- Test methodology

## Security Considerations
- [ ] No security impact
- [ ] Security improvement
- [ ] Potential security concern (please explain)

## Checklist
Please check all that apply:

- [ ] My code follows the project's coding standards
- [ ] I have performed a self-review of my code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings or errors
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published

## Additional Notes
Add any additional notes, concerns, or context for reviewers.

## Screenshots
If applicable, add screenshots to help explain your changes.

---

### For Maintainers
- [ ] Squash commits before merge
- [ ] Update version in relevant files
- [ ] Add appropriate labels
- [ ] Update CHANGELOG.md
- [ ] Update documentation if needed