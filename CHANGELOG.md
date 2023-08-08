# Changelog

## [1.9.2](https://github.com/cloudquery/plugin-pb-go/compare/v1.9.1...v1.9.2) (2023-08-08)


### Bug Fixes

* **deps:** Update github.com/apache/arrow/go/v13 digest to 112f949 ([#74](https://github.com/cloudquery/plugin-pb-go/issues/74)) ([a204890](https://github.com/cloudquery/plugin-pb-go/commit/a204890b005e5b1427228d42edcffb1671932a48))
* **deps:** Update golang.org/x/exp digest to b0cb94b ([#75](https://github.com/cloudquery/plugin-pb-go/issues/75)) ([2a62da5](https://github.com/cloudquery/plugin-pb-go/commit/2a62da512d6e49b7c9905d74a483fac5bd2f3fed))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to e0aa005 ([#77](https://github.com/cloudquery/plugin-pb-go/issues/77)) ([b273d9b](https://github.com/cloudquery/plugin-pb-go/commit/b273d9b0bf966ebaa8add42fb4747e6631e97fe6))
* **deps:** Update module github.com/google/flatbuffers to v23.5.26+incompatible ([#81](https://github.com/cloudquery/plugin-pb-go/issues/81)) ([4da6f65](https://github.com/cloudquery/plugin-pb-go/commit/4da6f6510e227539ce867e26754ac6532771358e))
* **deps:** Update module github.com/klauspost/compress to v1.16.7 ([#82](https://github.com/cloudquery/plugin-pb-go/issues/82)) ([12e9d67](https://github.com/cloudquery/plugin-pb-go/commit/12e9d670b6260fc6e15f3dcd5bfc6528b7c6cc05))
* **deps:** Update module github.com/mattn/go-runewidth to v0.0.15 ([#78](https://github.com/cloudquery/plugin-pb-go/issues/78)) ([f25ab8e](https://github.com/cloudquery/plugin-pb-go/commit/f25ab8e1f03ec4f789e15cebc913df34fcb82764))
* **deps:** Update module github.com/pierrec/lz4/v4 to v4.1.18 ([#79](https://github.com/cloudquery/plugin-pb-go/issues/79)) ([d94768e](https://github.com/cloudquery/plugin-pb-go/commit/d94768e5de70ce4ee771a8bb79381056711e8112))
* Generate Go Code from `plugin-pb` ([#83](https://github.com/cloudquery/plugin-pb-go/issues/83)) ([f898aa9](https://github.com/cloudquery/plugin-pb-go/commit/f898aa945651dfd646f3bb2736c6476a09d1c632))
* Generate Go Code from `plugin-pb` ([#84](https://github.com/cloudquery/plugin-pb-go/issues/84)) ([d0c9fa0](https://github.com/cloudquery/plugin-pb-go/commit/d0c9fa0087d1e37ef80b5adb5bc81ae76597986c))

## [1.9.1](https://github.com/cloudquery/plugin-pb-go/compare/v1.9.0...v1.9.1) (2023-07-28)


### Bug Fixes

* Wait for gRPC TCP connection in Docker registry ([#72](https://github.com/cloudquery/plugin-pb-go/issues/72)) ([98f3334](https://github.com/cloudquery/plugin-pb-go/commit/98f3334b7c1569989363982e4b20ac4d03e828fc))

## [1.9.0](https://github.com/cloudquery/plugin-pb-go/compare/v1.8.0...v1.9.0) (2023-07-28)


### Features

* Docker registry ([#70](https://github.com/cloudquery/plugin-pb-go/issues/70)) ([5bc1f60](https://github.com/cloudquery/plugin-pb-go/commit/5bc1f603676f6c49f8e504696e73341b93f670d4))

## [1.8.0](https://github.com/cloudquery/plugin-pb-go/compare/v1.7.1...v1.8.0) (2023-07-17)


### Features

* Add otel endpoint options for managed plugin ([#69](https://github.com/cloudquery/plugin-pb-go/issues/69)) ([4bb26b4](https://github.com/cloudquery/plugin-pb-go/commit/4bb26b4ff8824b63ac9c5f117ef11cd7f59b1706))


### Bug Fixes

* Generate Go Code from `plugin-pb` ([#67](https://github.com/cloudquery/plugin-pb-go/issues/67)) ([21dacdf](https://github.com/cloudquery/plugin-pb-go/commit/21dacdf120624ffdb2c837f4721f3d5a8d086e0f))

## [1.7.1](https://github.com/cloudquery/plugin-pb-go/compare/v1.7.0...v1.7.1) (2023-07-14)


### Bug Fixes

* Add `NoConnection` to init request ([#62](https://github.com/cloudquery/plugin-pb-go/issues/62)) ([80494f7](https://github.com/cloudquery/plugin-pb-go/commit/80494f7d37096db916aef6cd7148f560613335a9))

## [1.7.0](https://github.com/cloudquery/plugin-pb-go/compare/v1.6.0...v1.7.0) (2023-07-13)


### Features

* Add local path validation ([#60](https://github.com/cloudquery/plugin-pb-go/issues/60)) ([afa10cc](https://github.com/cloudquery/plugin-pb-go/commit/afa10ccd1d72b484174bc00668d35a19bcd56cb2))

## [1.6.0](https://github.com/cloudquery/plugin-pb-go/compare/v1.5.1...v1.6.0) (2023-07-03)


### Features

* Expose ConnectionString method in Client ([#57](https://github.com/cloudquery/plugin-pb-go/issues/57)) ([45b234e](https://github.com/cloudquery/plugin-pb-go/commit/45b234ee70ef717c91a247eb2bc7a4704db944df))

## [1.5.1](https://github.com/cloudquery/plugin-pb-go/compare/v1.5.0...v1.5.1) (2023-07-03)


### Bug Fixes

* Add backend options to Sync in V3 proto ([#56](https://github.com/cloudquery/plugin-pb-go/issues/56)) ([6182867](https://github.com/cloudquery/plugin-pb-go/commit/6182867ce3f62bcf110f8554759508bb7590b3fd))
* **deps:** Update github.com/apache/arrow/go/v13 digest to 5a06b2e ([#48](https://github.com/cloudquery/plugin-pb-go/issues/48)) ([cad785a](https://github.com/cloudquery/plugin-pb-go/commit/cad785afd36cc88abf34212e0e42a05266a8ef38))
* **deps:** Update golang.org/x/exp digest to 97b1e66 ([#49](https://github.com/cloudquery/plugin-pb-go/issues/49)) ([c9fc4d3](https://github.com/cloudquery/plugin-pb-go/commit/c9fc4d31522db282317afef441d6cc18b7222739))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to 9506855 ([#51](https://github.com/cloudquery/plugin-pb-go/issues/51)) ([da2709e](https://github.com/cloudquery/plugin-pb-go/commit/da2709effe871bacfb4ce0c2a656f3333b8d6783))
* **deps:** Update module github.com/goccy/go-json to v0.10.2 ([#52](https://github.com/cloudquery/plugin-pb-go/issues/52)) ([22857ff](https://github.com/cloudquery/plugin-pb-go/commit/22857ff2ddf478689971c40003bc35136e7fd0c4))
* **deps:** Update module github.com/klauspost/cpuid/v2 to v2.2.5 ([#53](https://github.com/cloudquery/plugin-pb-go/issues/53)) ([c1d951c](https://github.com/cloudquery/plugin-pb-go/commit/c1d951c08d904cef73bdbcaa7379ee44ba372b2b))
* **deps:** Update module github.com/mattn/go-isatty to v0.0.19 ([#54](https://github.com/cloudquery/plugin-pb-go/issues/54)) ([3406276](https://github.com/cloudquery/plugin-pb-go/commit/3406276e076bd62a3d25cd710caaac6ac13e3502))

## [1.5.0](https://github.com/cloudquery/plugin-pb-go/compare/v1.4.0...v1.5.0) (2023-06-30)


### Features

* Expose Read via gRPC ([#47](https://github.com/cloudquery/plugin-pb-go/issues/47)) ([5335b6d](https://github.com/cloudquery/plugin-pb-go/commit/5335b6dc66ff9013b8f600ac3cb2c76073fda7cb))


### Bug Fixes

* Add table_name to DeleteStale in proto v3 ([#45](https://github.com/cloudquery/plugin-pb-go/issues/45)) ([23eeffc](https://github.com/cloudquery/plugin-pb-go/commit/23eeffc90e0e8f8832fba7dd37cabd13f8a86974))

## [1.4.0](https://github.com/cloudquery/plugin-pb-go/compare/v1.3.6...v1.4.0) (2023-06-27)


### Features

* Add migrate_force to Write.InsertMessage ([#42](https://github.com/cloudquery/plugin-pb-go/issues/42)) ([600815d](https://github.com/cloudquery/plugin-pb-go/commit/600815dcc9faef6518a8ab6cbdba476a0ac6a483))

## [1.3.6](https://github.com/cloudquery/plugin-pb-go/compare/v1.3.5...v1.3.6) (2023-06-27)


### Bug Fixes

* Split Sync and Write messages to it's own proto messages ([#40](https://github.com/cloudquery/plugin-pb-go/issues/40)) ([1bd6271](https://github.com/cloudquery/plugin-pb-go/commit/1bd62719f0eac5d6f58e10abf6b48566e5ee3352))

## [1.3.5](https://github.com/cloudquery/plugin-pb-go/compare/v1.3.4...v1.3.5) (2023-06-27)


### Bug Fixes

* Remove migrate_force from plugin v3 ([#37](https://github.com/cloudquery/plugin-pb-go/issues/37)) ([6e1cf13](https://github.com/cloudquery/plugin-pb-go/commit/6e1cf13c8359d9387730173e3fa3c8fb5de8a4be))

## [1.3.4](https://github.com/cloudquery/plugin-pb-go/compare/v1.3.3...v1.3.4) (2023-06-26)


### Bug Fixes

* Regenerate V3, remove backend from proto ([#35](https://github.com/cloudquery/plugin-pb-go/issues/35)) ([78ae019](https://github.com/cloudquery/plugin-pb-go/commit/78ae019b01322dd8ab4f48daa5a6e00d02f8a2cf))

## [1.3.3](https://github.com/cloudquery/plugin-pb-go/compare/v1.3.2...v1.3.3) (2023-06-24)


### Bug Fixes

* Add record enc/dec to destv1 sourcev2 ([#33](https://github.com/cloudquery/plugin-pb-go/issues/33)) ([40797be](https://github.com/cloudquery/plugin-pb-go/commit/40797be0bb62422984845597fbd984b877c76032))

## [1.3.2](https://github.com/cloudquery/plugin-pb-go/compare/v1.3.1...v1.3.2) (2023-06-23)


### Bug Fixes

* Add Schemas encoding/decoding to plugin v3 proto ([#30](https://github.com/cloudquery/plugin-pb-go/issues/30)) ([a549e89](https://github.com/cloudquery/plugin-pb-go/commit/a549e89cb7b34e72db9d9018bf6be47b12182f3d))

## [1.3.1](https://github.com/cloudquery/plugin-pb-go/compare/v1.3.0...v1.3.1) (2023-06-23)


### Bug Fixes

* Update schema encoding/deocoding in v2 ([#28](https://github.com/cloudquery/plugin-pb-go/issues/28)) ([6678004](https://github.com/cloudquery/plugin-pb-go/commit/66780042358299e25d0e2d30f4ecd49c15766f77))

## [1.3.0](https://github.com/cloudquery/plugin-pb-go/compare/v1.2.1...v1.3.0) (2023-06-23)


### Features

* Add arrow schema encoding to v2 ([#26](https://github.com/cloudquery/plugin-pb-go/issues/26)) ([a7399f5](https://github.com/cloudquery/plugin-pb-go/commit/a7399f57a6f612f579321b0dedf11e425f1e6a32))

## [1.2.1](https://github.com/cloudquery/plugin-pb-go/compare/v1.2.0...v1.2.1) (2023-06-23)


### Bug Fixes

* Discovery V1 regen ([#24](https://github.com/cloudquery/plugin-pb-go/issues/24)) ([5c5dd27](https://github.com/cloudquery/plugin-pb-go/commit/5c5dd27d950f8ef2d528cdfe512d2ff51346e3d8))

## [1.2.0](https://github.com/cloudquery/plugin-pb-go/compare/v1.1.0...v1.2.0) (2023-06-23)


### Features

* Add Plugin Proto V3 ([#21](https://github.com/cloudquery/plugin-pb-go/issues/21)) ([50ec9d9](https://github.com/cloudquery/plugin-pb-go/commit/50ec9d90942e74677e39e8379cba6631cde40e04))

## [1.1.0](https://github.com/cloudquery/plugin-pb-go/compare/v1.0.9...v1.1.0) (2023-06-19)


### Features

* Add managedplugin  ([#16](https://github.com/cloudquery/plugin-pb-go/issues/16)) ([afb3415](https://github.com/cloudquery/plugin-pb-go/commit/afb3415accd4932862cf6df23660dd242164dd6e))

## [1.0.9](https://github.com/cloudquery/plugin-pb-go/compare/v1.0.8...v1.0.9) (2023-06-06)


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to 2e198f4 ([#6](https://github.com/cloudquery/plugin-pb-go/issues/6)) ([bbf4975](https://github.com/cloudquery/plugin-pb-go/commit/bbf4975f895c4a930962ccd00bcdfcc33154715b))
* **deps:** Update google.golang.org/genproto digest to e85fd2c ([#7](https://github.com/cloudquery/plugin-pb-go/issues/7)) ([ead0d7f](https://github.com/cloudquery/plugin-pb-go/commit/ead0d7f4f142f95c9c8ac52da0dced22ddccae61))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to e85fd2c ([#11](https://github.com/cloudquery/plugin-pb-go/issues/11)) ([e21aa76](https://github.com/cloudquery/plugin-pb-go/commit/e21aa7631b9ea7b5824d625b923cfb088c6f6108))
* **deps:** Update module github.com/davecgh/go-spew to v1.1.1 ([#9](https://github.com/cloudquery/plugin-pb-go/issues/9)) ([a6a34f4](https://github.com/cloudquery/plugin-pb-go/commit/a6a34f4d7e4988a649a26b91115af6a4eb7860aa))
* **deps:** Update module github.com/stretchr/testify to v1.8.4 ([#12](https://github.com/cloudquery/plugin-pb-go/issues/12)) ([dba4977](https://github.com/cloudquery/plugin-pb-go/commit/dba497785ca2b781c24d7d8120488502eb5b24c4))
* Embedded content handling with newlines ([#13](https://github.com/cloudquery/plugin-pb-go/issues/13)) ([ceb6046](https://github.com/cloudquery/plugin-pb-go/commit/ceb6046ab407c14df991fb1ee2caf494a0aa278a))
* SpecReader should escape external JSON content from files and environment variables ([#4](https://github.com/cloudquery/plugin-pb-go/issues/4)) ([54b172f](https://github.com/cloudquery/plugin-pb-go/commit/54b172f13b19b2ee59009098679f45fae67f28a3))
