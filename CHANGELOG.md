## [Unreleased]

### Added

* Add `archived_query_as_deleted` provider argument to treat an archived query as deleted, so `redash_query` drops it from the state and the next plan creates it again. Redash has no API that deletes a query. [pull#236](https://github.com/winebarrel/terraform-provider-redash/pull/236)

### Fixed

* Remove `redash_data_source`, `redash_group_data_sources`, `redash_group_data_source`, `redash_user` and `redash_alert_destination` from the state when they are deleted outside Terraform, instead of failing on refresh. For the two grant resources this means their group was deleted. [pull#234](https://github.com/winebarrel/terraform-provider-redash/pull/234)

* Fix `redash_query` keeping a schedule in the state after it is removed outside Terraform, so the plan showed no diff. [pull#235](https://github.com/winebarrel/terraform-provider-redash/pull/235)

## [1.15.0] - 2026-10-26

### Fixed

* Fix `redash_group_data_sources` failing to plan with `duplicate data_source_id 0` when its `data_source` blocks reference two or more data sources created in the same apply with different `view_only` values. [pull#230](https://github.com/winebarrel/terraform-provider-redash/pull/230)

* Fix `redash_group_data_sources` and `redash_group_data_source` failing to remove a grant when the data source is deleted in the same apply. A grant that is already gone now counts as removed. [pull#230](https://github.com/winebarrel/terraform-provider-redash/pull/230)

## [1.14.0] - 2026-10-22

### Changed

* **BREAKING** Replace `in_default_group` on `redash_data_source` with `remove_from_default_group_on_create`, which only works on create. [pull#229](https://github.com/winebarrel/terraform-provider-redash/pull/229)

## [1.13.0] - 2026-10-21

### Added

* Add `in_default_group` argument to `redash_data_source` to keep the data source out of the `default` group, which Redash adds every new data source to. [pull#215](https://github.com/winebarrel/terraform-provider-redash/pull/215) [pull#223](https://github.com/winebarrel/terraform-provider-redash/pull/223)

### Changed

* Read `redash_group_member` from Redash, so a membership removed outside Terraform shows up as a diff and importing a non-existent membership fails. [pull#216](https://github.com/winebarrel/terraform-provider-redash/pull/216) [pull#218](https://github.com/winebarrel/terraform-provider-redash/pull/218) [pull#222](https://github.com/winebarrel/terraform-provider-redash/pull/222)
* Read `redash_alert_subscription` from Redash, so a subscription removed outside Terraform shows up as a diff. [pull#221](https://github.com/winebarrel/terraform-provider-redash/pull/221)
* Remove `redash_group` and `redash_alert` from the state when they are deleted outside Terraform, instead of failing on refresh. [pull#219](https://github.com/winebarrel/terraform-provider-redash/pull/219) [pull#220](https://github.com/winebarrel/terraform-provider-redash/pull/220)
* Return an error instead of ignoring malformed IDs on `redash_group_member` and `redash_alert_subscription` delete. [pull#222](https://github.com/winebarrel/terraform-provider-redash/pull/222)
* Add descriptions to `redash_group_member` and `redash_alert_subscription` arguments. [pull#222](https://github.com/winebarrel/terraform-provider-redash/pull/222)

### Fixed

* Fix a panic on `redash_alert_subscription` import when the alert has a subscription without a destination. [pull#221](https://github.com/winebarrel/terraform-provider-redash/pull/221)

### Security

* Build with Go 1.27.3 to fix CVE-2026-111943 ([GO-2026-6809](https://pkg.go.dev/vuln/GO-2026-6809)) in `crypto/x509`. [pull#217](https://github.com/winebarrel/terraform-provider-redash/pull/217)

## [1.12.0] - 2026-10-05

### Added

* Add write-only `options_wo` and `options_wo_version` arguments to `redash_data_source` so that data source options (e.g. credentials) are not stored in the Terraform state. Requires Terraform 1.11 or later. [pull#210](https://github.com/winebarrel/terraform-provider-redash/pull/210) [pull#211](https://github.com/winebarrel/terraform-provider-redash/pull/211) [pull#212](https://github.com/winebarrel/terraform-provider-redash/pull/212)

## [1.11.0] - 2026-10-03

### Added

* Add `redash_group_data_sources` resource to manage all the data sources granted to a group. [pull#205](https://github.com/winebarrel/terraform-provider-redash/pull/205) [pull#206](https://github.com/winebarrel/terraform-provider-redash/pull/206) [pull#209](https://github.com/winebarrel/terraform-provider-redash/pull/209)

## [1.10.0] - 2026-10-02

### Added

* Add `redash_group_data_sources` data source to list the data sources granted to a group. [pull#202](https://github.com/winebarrel/terraform-provider-redash/pull/202)

## [1.9.0] - 2026-10-02

### Changed

* Update redash-go to [v2.12.0](https://github.com/winebarrel/redash-go/blob/main/CHANGELOG.md#2120---2026-10-02).
* Update `redash_group` name in place instead of recreating the group. [pull#201](https://github.com/winebarrel/terraform-provider-redash/pull/201)

### Fixed

* Fix `redash_group` being replaced after import. [pull#196](https://github.com/winebarrel/terraform-provider-redash/pull/196)
* Fix `view_only` not being set on `redash_group_data_source` import. [pull#197](https://github.com/winebarrel/terraform-provider-redash/pull/197)

## [1.8.0] - 2026-10-01

### Changed

* Add `http_headers` provider argument to send extra HTTP headers on Redash API requests. [pull#190](https://github.com/winebarrel/terraform-provider-redash/pull/190)

## [1.7.0] - 2026-09-17

### Changed

* Update Go to 1.27. [pull#187](https://github.com/winebarrel/terraform-provider-redash/pull/187)
* Return an error instead of panicking on malformed import IDs. [pull#187](https://github.com/winebarrel/terraform-provider-redash/pull/187)

### Fixed

* Fix pagination in `redash_query` data source. [pull#186](https://github.com/winebarrel/terraform-provider-redash/pull/186)

## [1.6.0] - 2026-09-17

### Changed

* Update redash-go to [v2.11.0](https://github.com/winebarrel/redash-go/blob/main/CHANGELOG.md#2110---2026-09-17).
* Add `selector` to `redash_alert` options. [pull#179](https://github.com/winebarrel/terraform-provider-redash/pull/179)

## [1.5.0] - 2026-09-11

### Fixed

* Fix permanent drift on `data_source` resource options with secrets. [pull#178](https://github.com/winebarrel/terraform-provider-redash/pull/178)

## [1.4.0] - 2026-04-29

### Changed

* Update redash-go to [v2.10.0](https://github.com/winebarrel/redash-go/blob/main/CHANGELOG.md#2100---2026-04-29).
* Don't send schedule with interval 0 for queries without schedule block. [pull#149](https://github.com/winebarrel/terraform-provider-redash/pull/149)
* Clear schedule on update when schedule block is removed. [pull#152](https://github.com/winebarrel/terraform-provider-redash/pull/152)

## [1.3.0] - 2026-02-25

### Changed

* Update redash-go to [v2.9.0](https://github.com/winebarrel/redash-go/blob/main/CHANGELOG.md#290---2026-02-25).
* Change alert value schema from TypeInt to TypeFloat. [pull#143](https://github.com/winebarrel/terraform-provider-redash/pull/143)

## [1.2.0] - 2025-03-10

### Added

* Support parameterized query.

## [1.1.0] - 2025-03-09

### Added

* Add `publish` attr to `query` resource/data source.

# Changelog

## [1.0.0] - 2025-01-21

### Changed

* Update redash-go to v2.

## [0.4.1] - 2023-04-26

### Changed

* Update redash-go.
* Update query tags update processing.

## [0.4.0] - 2023-04-26

### Added

* Support `redash_alert_subscription` import.

### Changed

* Check import ID.

## [0.3.2] - 2023-04-26

### Changed

* Remove resource/data source description.

## [0.3.1] - 2023-04-25

### Changed

* Update document.

## [0.3.0] - 2023-04-25

### Added

* Support query tags.

## [0.2.0] - 2023-04-25

### Added

* Add `redash_user` resource.
* Add debug output.

### Changed

* Ignore disabled users.

## [0.1.0] - 2023-04-25

### Added

- First release.

<!-- cf. https://keepachangelog.com/ -->
