# Change Log
All releases of the BOSH CPI for Google Cloud Platform will be documented in
this file. This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]
### Added
- This changelog

## [24.1.0] - 2015-12-03

### Changed
- Instance tags can be specified in any `cloud_properties` section of a BOSH manifest

### Removed
- The dummy BOSH release is no longer part of the CI pipeline

### Fixed
- Integration tests will use the CI pipeline stemcell rather than requiring an existing stemcell in a project

[Unreleased]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24.1.0...HEAD
[24.1.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24...v24.1.0
