# Releasing Guidelines for Cucumber Godog

This document provides guidelines for releasing new versions of Cucumber Godog. Follow these steps to ensure a smooth and consistent release process.

## Versioning

Cucumber Godog follows [Semantic Versioning]. Version numbers are in the format `MAJOR.MINOR.PATCH`.

### Current (for v0.MINOR.PATCH)

- **MINOR**: Incompatible API changes.
- **PATCH**: Backward-compatible new features and bug fixes.

### After v1.X.X release

- **MAJOR**: Incompatible API changes.
- **MINOR**: Backward-compatible new features.
- **PATCH**: Backward-compatible bug fixes.

## Release Process

1. **Update Changelog:**
    - Open `CHANGELOG.md` and add an entry for the upcoming release formatting according to the principles of [Keep A CHANGELOG].
    - Include details about new features, enhancements, and bug fixes.

2. **Run Tests:**
    - Run the test suite to ensure all existing features are working as expected.

3. **Manual Testing for Backwards Compatibility:**
    - Manually test the new release with external libraries that depend on Cucumber Godog.
    - Look for any potential backwards compatibility issues, especially with widely-used libraries.
    - Address any identified issues before proceeding.

4. **Create Release on GitHub:**
    - Go to the [Releases] page on GitHub.
    - Click on "Draft a new release."
    - Tag version should be set to the new tag vMAJOR.MINOR.PATCH
    - Title the release using the version number (e.g., "vMAJOR.MINOR.PATCH").
    - Click 'Generate release notes'

5. **Publish Release:**
    - Click "Publish release" to make the release public.

6. **Announce the Release:**
    - Make an announcement on relevant communication channels (e.g., [community Discord]) about the new release.

## Additional Considerations

- **Documentation:**
    - Update the project documentation on the [website], if applicable.

- **Deprecation Notices:**
    - If any features are deprecated, clearly document them in the release notes and provide guidance on migration.

- **Compatibility:**
    - Clearly state any compatibility requirements or changes in the release notes.

- **Feedback:**
    - Encourage users to provide feedback and report any issues with the new release.

Following these guidelines, including manual testing with external libraries, will help ensure a thorough release process for Cucumber Godog, allowing detection and resolution of potential backwards compatibility issues before tagging the release.

[community Discord]: https://cucumber.io/community#discord
[website]: https://cucumber.github.io/godog/
[Releases]: https://github.com/cucumber/godog/releases
[Semantic Versioning]: http://semver.org
[Keep A CHANGELOG]: http://keepachangelog.com