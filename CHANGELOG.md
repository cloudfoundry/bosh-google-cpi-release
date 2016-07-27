# Change Log
All releases of the BOSH CPI for Google Cloud Platform will be documented in
this file. This project adheres to [Semantic Versioning](http://semver.org/).

## [24.3.0] - 2016-07-25

### Added
- A `resource_pool`'s manifest can now specify `ephemeral_external_ip` and
  `ip_forwarding` properties, overriding those same properties in a
  manifest's `networks` section.

## [24.2.0] - 2016-07-25

### Added
- This changelog

### Changed
- 3262.4 stemcell

### Fixed
- All tests now use light stemcells

## [24.1.0] - 2016-07-25

### Changed
- Instance tags can be specified in any `cloud_properties` section of a BOSH manifest

### Removed
- The dummy BOSH release is no longer part of the CI pipeline

### Fixed
- Integration tests will use the CI pipeline stemcell rather than requiring an existing stemcell in a project

[24.3.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24.2.0...v24.3.0
[24.2.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24.1.0...v24.2.0
[24.1.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24...v24.1.0
