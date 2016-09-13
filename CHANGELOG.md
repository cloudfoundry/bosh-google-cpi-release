# Change Log
All releases of the BOSH CPI for Google Cloud Platform will be documented in
this file. This project adheres to [Semantic Versioning](http://semver.org/).

## [25.2.0] - 2016-08-18

### Changed
- Any metadata provided by bosh in the `set_vm_metadata` action will also be 
  propagated to the VM as [labels](https://cloud.google.com/compute/docs/label-or-tag-resources),
  allowing sorting and filter in the web console based on job, deployment, etc.

## [25.1.0] - 2016-08-18

### Added
- The `service_account` cloud-config property may now use the e-mail address
  of a custom service account.

## [25.0.0] - 2016-07-25

### Changed
- The `default_zone` config property (in the `google` section of a manifest)
  is no longer supported. The `zone` property must be explicitly set in the
  `cloud_properties` section of `resource_pools` (or `azs` for a cloud-config
  director.)

## [24.4.0] - 2016-07-25

### Fixed
- An explicit region is used to locate a subnet, allowing subnets with the same
  name to be differentiated.

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

[25.2.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v25.1.0...v25.2.0
[25.1.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v25.0.0...v25.1.0
[25.0.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24.4.0...v25.0.0
[24.4.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24.3.0...v24.4.0
[24.3.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24.2.0...v24.3.0
[24.2.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24.1.0...v24.2.0
[24.1.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24...v24.1.0
