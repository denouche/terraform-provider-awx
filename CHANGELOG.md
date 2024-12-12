# [0.28.0](https://github.com/denouche/terraform-provider-awx/compare/v0.27.0...v0.28.0) (2024-12-10)


### Features

* link instance group with organization ([#85](https://github.com/denouche/terraform-provider-awx/issues/85)) ([21dbb14](https://github.com/denouche/terraform-provider-awx/commit/21dbb1492ff66a07002d1898eb3d7588158b3ee3))

# [0.27.0](https://github.com/denouche/terraform-provider-awx/compare/v0.26.1...v0.27.0) (2024-05-31)


### Features

* **awx_job_template_launch:** add job template launch timeout ([#72](https://github.com/denouche/terraform-provider-awx/issues/72)) ([e23618c](https://github.com/denouche/terraform-provider-awx/commit/e23618c7c71270389fe93a77314ad8cf44ea1bfb))

## [0.26.1](https://github.com/denouche/terraform-provider-awx/compare/v0.26.0...v0.26.1) (2024-05-31)


### Bug Fixes

* organization default env type ([af1c640](https://github.com/denouche/terraform-provider-awx/commit/af1c6402f482c894a28e113abcfaf145bde7783b))
* organization default env var ([efed961](https://github.com/denouche/terraform-provider-awx/commit/efed961273f09b7b87dc1fb0eb657b6179647cbe))

# [0.26.0](https://github.com/denouche/terraform-provider-awx/compare/v0.25.2...v0.26.0) (2024-03-19)


### Features

* fix schedule definition and allow job template schedule import ([#56](https://github.com/denouche/terraform-provider-awx/issues/56)) ([219d7c3](https://github.com/denouche/terraform-provider-awx/commit/219d7c3acf83a875f9b008c8af7ebf8da25f7b42))

## [0.25.2](https://github.com/denouche/terraform-provider-awx/compare/v0.25.1...v0.25.2) (2024-03-19)


### Bug Fixes

* **resource_organization:** undefined attribut project_id [#63](https://github.com/denouche/terraform-provider-awx/issues/63) ([#64](https://github.com/denouche/terraform-provider-awx/issues/64)) ([6276aa1](https://github.com/denouche/terraform-provider-awx/commit/6276aa15ec5ed06a300b90807acdd912747e04bd))

## [0.25.1](https://github.com/denouche/terraform-provider-awx/compare/v0.25.0...v0.25.1) (2024-03-19)


### Bug Fixes

* **resource_job_template_launch:** rename inventory to inventory_id and add support for yaml variables ([#59](https://github.com/denouche/terraform-provider-awx/issues/59)) ([6c4f141](https://github.com/denouche/terraform-provider-awx/commit/6c4f141f52338e51a89ab42c0c4858fc006052df))

# [0.25.0](https://github.com/denouche/terraform-provider-awx/compare/v0.24.4...v0.25.0) (2024-03-19)


### Features

* add description and name fields to datasource awx_credential(s) ([#55](https://github.com/denouche/terraform-provider-awx/issues/55)) ([02715cc](https://github.com/denouche/terraform-provider-awx/commit/02715cce3239d0f7540ffb058b1c02a368ad7a26))

## [0.24.4](https://github.com/denouche/terraform-provider-awx/compare/v0.24.3...v0.24.4) (2024-03-19)

## [0.24.3](https://github.com/denouche/terraform-provider-awx/compare/v0.24.2...v0.24.3) (2023-11-23)


### Bug Fixes

* allow empty credential_id for inventory_source creation ([b960b53](https://github.com/denouche/terraform-provider-awx/commit/b960b538a60d7d28c566eac3c887985d58bf4d6a)), closes [#19](https://github.com/denouche/terraform-provider-awx/issues/19)
* allow empty source_project_id for inventory_source creation ([#29](https://github.com/denouche/terraform-provider-awx/issues/29)) ([4d8d1d8](https://github.com/denouche/terraform-provider-awx/commit/4d8d1d8de62b50e0f9614175f993e10a4c92c609))

## [0.24.2](https://github.com/denouche/terraform-provider-awx/compare/v0.24.1...v0.24.2) (2023-11-23)


### Bug Fixes

* fix compilation issue in resource project ([45771cd](https://github.com/denouche/terraform-provider-awx/commit/45771cd8ec7b9409c073fb7110105788d546e586))

## [0.24.1](https://github.com/denouche/terraform-provider-awx/compare/v0.24.0...v0.24.1) (2023-11-23)


### Bug Fixes

* fix resource project ([6519171](https://github.com/denouche/terraform-provider-awx/commit/651917165ca66be157b428558a4e1932b294e01c))

# [0.24.0](https://github.com/denouche/terraform-provider-awx/compare/v0.23.0...v0.24.0) (2023-11-23)


### Features

* add Organizations GalaxyCredentials, resource credential Ansible Galaxy, user resource, organization role, Gitlab credential, resource settings, instance group support ([5a30c50](https://github.com/denouche/terraform-provider-awx/commit/5a30c505fb9d2e4fdc6b72b080e16effaf47d1d4))

# [0.23.0](https://github.com/denouche/terraform-provider-awx/compare/v0.22.6...v0.23.0) (2023-11-23)


### Features

* add allow_override for awx_project. ([#30](https://github.com/denouche/terraform-provider-awx/issues/30)) ([285bb45](https://github.com/denouche/terraform-provider-awx/commit/285bb45ef85e9988434fb95d1b711fd411247217))

## [0.22.6](https://github.com/denouche/terraform-provider-awx/compare/v0.22.5...v0.22.6) (2023-11-23)

## [0.22.5](https://github.com/denouche/terraform-provider-awx/compare/v0.22.4...v0.22.5) (2023-11-23)


### Bug Fixes

* segmentation fault when datasource does not exist ([fa53eac](https://github.com/denouche/terraform-provider-awx/commit/fa53eac9aa1ec3be4ceedda15e6b6a2173c9ed18))

## [0.22.4](https://github.com/denouche/terraform-provider-awx/compare/v0.22.3...v0.22.4) (2023-11-23)

## [0.22.3](https://github.com/denouche/terraform-provider-awx/compare/v0.22.2...v0.22.3) (2023-11-23)

## [0.22.2](https://github.com/denouche/terraform-provider-awx/compare/v0.22.1...v0.22.2) (2023-11-23)

## [0.22.1](https://github.com/denouche/terraform-provider-awx/compare/v0.22.0...v0.22.1) (2023-11-23)

# [0.22.0](https://github.com/denouche/terraform-provider-awx/compare/v0.21.0...v0.22.0) (2023-11-23)


### Features

* **resource_job_template_launch:** add extra_vars option ([3e4ab25](https://github.com/denouche/terraform-provider-awx/commit/3e4ab257f760fd5aca4f206ff1d00d2467f4652f))
* **resource_job_template_launch:** add more configuration options (limit, inventory, completion wait) ([abe120c](https://github.com/denouche/terraform-provider-awx/commit/abe120ce93029b02b9b34f29be4615e430a9d8c5))

# [0.21.0](https://github.com/denouche/terraform-provider-awx/compare/v0.20.0...v0.21.0) (2023-04-17)


### Features

* schedule extra data ([1d60ab8](https://github.com/denouche/terraform-provider-awx/commit/1d60ab86f97446abdca431020d0e4e4537096e04))

# [0.20.0](https://github.com/denouche/terraform-provider-awx/compare/v0.19.0...v0.20.0) (2023-04-13)


### Bug Fixes

* resource data description ([bd38f22](https://github.com/denouche/terraform-provider-awx/commit/bd38f22a2b5f1a0828998ccc555d1879588510b5))


### Features

* extra_data to workflow schedule schema ([df67648](https://github.com/denouche/terraform-provider-awx/commit/df6764890be3007f09f284e59f7bbef8eac2586c))

# [0.19.0](https://github.com/denouche/terraform-provider-awx/compare/v0.18.0...v0.19.0) (2022-11-14)


### Features

* Add organization role data source ([4bc4065](https://github.com/denouche/terraform-provider-awx/commit/4bc40653f96d92c47d0e0f5fb53d4172661491f4))
* Add setting resource ([b1b1a24](https://github.com/denouche/terraform-provider-awx/commit/b1b1a2403887599bd451e54094a48c7e728aa8da))
* fetch upstream ([8cc9cb0](https://github.com/denouche/terraform-provider-awx/commit/8cc9cb0f160b779e02f17fab10c03d4cb7ec54b9)), closes [#16](https://github.com/denouche/terraform-provider-awx/issues/16)

# [0.18.0](https://github.com/denouche/terraform-provider-awx/compare/v0.17.0...v0.18.0) (2022-11-14)


### Features

* add AWX token authentication ([#15](https://github.com/denouche/terraform-provider-awx/issues/15)) ([55b7d41](https://github.com/denouche/terraform-provider-awx/commit/55b7d41579f79d7fbb1aa61bd243405a81815748))

# [0.17.0](https://github.com/denouche/terraform-provider-awx/compare/v0.16.0...v0.17.0) (2022-10-31)


### Features

* adds the possibility to use source_id inside resource_inventory_source ([#20](https://github.com/denouche/terraform-provider-awx/issues/20)) ([6891c9e](https://github.com/denouche/terraform-provider-awx/commit/6891c9eb98b8ca916746af08c640520f07d29dda))

# [0.16.0](https://github.com/denouche/terraform-provider-awx/compare/v0.15.6...v0.16.0) (2022-10-25)


### Features

* upgrade goawx lib to 0.14.1 ([#22](https://github.com/denouche/terraform-provider-awx/issues/22)) ([3193f56](https://github.com/denouche/terraform-provider-awx/commit/3193f56a55ac96103f0b2a4f355f0cd723116f86))

## [0.15.6](https://github.com/denouche/terraform-provider-awx/compare/v0.15.5...v0.15.6) (2022-07-22)


### Bug Fixes

* when using insecure connection the PROXY_HTTPS env var was ignored ([#12](https://github.com/denouche/terraform-provider-awx/issues/12)) ([e457deb](https://github.com/denouche/terraform-provider-awx/commit/e457deb4644f82e4c4e3af27e07df7ba565cbbaa))

## [0.15.5](https://github.com/denouche/terraform-provider-awx/compare/v0.15.4...v0.15.5) (2022-07-22)


### Bug Fixes

* fix some func names after upgrading goawx dep ([999e70d](https://github.com/denouche/terraform-provider-awx/commit/999e70ddbdcdc3ca758b85e9c6a4eea3b3859689))

## [0.15.4](https://github.com/denouche/terraform-provider-awx/compare/v0.15.3...v0.15.4) (2022-07-22)

## [0.15.3](https://github.com/denouche/terraform-provider-awx/compare/v0.15.2...v0.15.3) (2022-07-19)


### Bug Fixes

* do not provide local_path for project if the scm_type is git ([#13](https://github.com/denouche/terraform-provider-awx/issues/13)) ([b4ab7dc](https://github.com/denouche/terraform-provider-awx/commit/b4ab7dc51306507bd71ef61b611782567bc0c0bb))

## [0.15.2](https://github.com/denouche/terraform-provider-awx/compare/v0.15.1...v0.15.2) (2022-07-01)


### Bug Fixes

* make a new release ([be91fb4](https://github.com/denouche/terraform-provider-awx/commit/be91fb4577e932ffee1019efb70620479d6089fd))

## [0.15.1](https://github.com/denouche/terraform-provider-awx/compare/v0.15.0...v0.15.1) (2022-07-01)


### Bug Fixes

* goawx version for always node type ([#7](https://github.com/denouche/terraform-provider-awx/issues/7)) ([bfe6ea8](https://github.com/denouche/terraform-provider-awx/commit/bfe6ea8d2245836a5b2584b4d471ca911d1b4626))

# [0.15.0](https://github.com/denouche/terraform-provider-awx/compare/v0.14.0...v0.15.0) (2022-05-11)


### Features

* organizations data source ([#4](https://github.com/denouche/terraform-provider-awx/issues/4)) ([ad61e88](https://github.com/denouche/terraform-provider-awx/commit/ad61e88a638b94eda2c306a0d9f610d65508d17f))

# [0.14.0](https://github.com/denouche/terraform-provider-awx/compare/v0.13.1...v0.14.0) (2022-04-21)


### Features

* workflow job template notifications ([#3](https://github.com/denouche/terraform-provider-awx/issues/3)) ([00db915](https://github.com/denouche/terraform-provider-awx/commit/00db9157df52d9fb4431db6f53ac5aa8038bad44))

## [0.13.1](https://github.com/denouche/terraform-provider-awx/compare/v0.13.0...v0.13.1) (2022-04-20)


### Bug Fixes

* workflow job template & schedule inventory option default value ([#2](https://github.com/denouche/terraform-provider-awx/issues/2)) ([6869420](https://github.com/denouche/terraform-provider-awx/commit/6869420d6b87a70922c915d1012ebd15156a277a))

# [0.13.0](https://github.com/denouche/terraform-provider-awx/compare/v0.12.3...v0.13.0) (2022-04-20)


### Features

* support execution environments ([#1](https://github.com/denouche/terraform-provider-awx/issues/1)) ([0791c09](https://github.com/denouche/terraform-provider-awx/commit/0791c09cb85783e7433f8e4ea80cfa9d7911af32))

## [0.12.3](https://github.com/denouche/terraform-provider-awx/compare/v0.12.2...v0.12.3) (2022-04-19)


### Bug Fixes

* publish for all os and arch ([7a3cd45](https://github.com/denouche/terraform-provider-awx/commit/7a3cd4552b44246377a00a185dbde48b45ce07dc))

## [0.12.2](https://github.com/denouche/terraform-provider-awx/compare/v0.12.1...v0.12.2) (2022-01-05)


### Bug Fixes

* upgrade goawx dep ([ba2ea50](https://github.com/denouche/terraform-provider-awx/commit/ba2ea509f164f7dad4f5477d6d58a40a798c0022))

## [0.12.1](https://github.com/denouche/terraform-provider-awx/compare/v0.12.0...v0.12.1) (2022-01-05)


### Bug Fixes

* upgrade goawx dep ([50447a2](https://github.com/denouche/terraform-provider-awx/commit/50447a2ebf2a0fb2862f2749a6aaa7ec58fed0e7))

# [0.12.0](https://github.com/denouche/terraform-provider-awx/compare/v0.11.4...v0.12.0) (2022-01-05)


### Features

* add resources awx_job_template_notification_template_success awx_job_template_notification_template_error awx_job_template_notification_template_started ([24b69c5](https://github.com/denouche/terraform-provider-awx/commit/24b69c5ded4c0fbba366637c0e423e0fc07679e6))

## [0.11.4](https://github.com/denouche/terraform-provider-awx/compare/v0.11.3...v0.11.4) (2021-12-24)


### Bug Fixes

* notification template notification configuration is a json ([09787ef](https://github.com/denouche/terraform-provider-awx/commit/09787ef93e745a0049970f5fcd134f5ab5a7f6f5))

## [0.11.3](https://github.com/denouche/terraform-provider-awx/compare/v0.11.2...v0.11.3) (2021-12-24)


### Bug Fixes

* notification_configuration is a string ([f10fb3b](https://github.com/denouche/terraform-provider-awx/commit/f10fb3ba03deca84d3169bc2eac0b01503c438f8))

## [0.11.2](https://github.com/denouche/terraform-provider-awx/compare/v0.11.1...v0.11.2) (2021-12-24)


### Bug Fixes

* some fixes on notification_template resource ([3cd1a59](https://github.com/denouche/terraform-provider-awx/commit/3cd1a592ad1c3baed7a237aa228645a90cb790cb))

## [0.11.1](https://github.com/denouche/terraform-provider-awx/compare/v0.11.0...v0.11.1) (2021-12-24)


### Bug Fixes

* notification_template schema ([4b28594](https://github.com/denouche/terraform-provider-awx/commit/4b2859405fc56bb7a09320f826862cbaa05a6d32))

# [0.11.0](https://github.com/denouche/terraform-provider-awx/compare/v0.10.7...v0.11.0) (2021-12-24)


### Features

* add notification_template resource ([9c5b488](https://github.com/denouche/terraform-provider-awx/commit/9c5b4885dfcd068b7dbac89567067c606b73fa6c))

## [0.10.7](https://github.com/denouche/terraform-provider-awx/compare/v0.10.6...v0.10.7) (2021-12-23)


### Bug Fixes

* add missing unified_job_template_id ([72d1deb](https://github.com/denouche/terraform-provider-awx/commit/72d1deb810d8618158bb48cea924959961495163))

## [0.10.6](https://github.com/denouche/terraform-provider-awx/compare/v0.10.5...v0.10.6) (2021-12-23)


### Bug Fixes

* add schedule inventory parameter support ([5e691ac](https://github.com/denouche/terraform-provider-awx/commit/5e691ac67f0e28337688928a96d6a3f1b0a7376a))

## [0.10.5](https://github.com/denouche/terraform-provider-awx/compare/v0.10.4...v0.10.5) (2021-12-23)

## [0.10.4](https://github.com/denouche/terraform-provider-awx/compare/v0.10.3...v0.10.4) (2021-12-23)

## [0.10.3](https://github.com/denouche/terraform-provider-awx/compare/v0.10.2...v0.10.3) (2021-12-23)
