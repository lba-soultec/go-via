# [1.1.0](https://github.com/lba-soultec/go-via/compare/v1.0.12...v1.1.0) (2025-08-13)


### Bug Fixes

* **ai:** resolve some issues ai did not cover and make it runnable. still a proposal ([afdc876](https://github.com/lba-soultec/go-via/commit/afdc87641338a1f3b188244ab0ea76b7f6ab6000))
* **api:** remove ai mistakes ([9a2a991](https://github.com/lba-soultec/go-via/commit/9a2a9918252cdae9694aac5e448aa508f6cb55b8))
* **debug:** remove debug ([085e33d](https://github.com/lba-soultec/go-via/commit/085e33d891961b2d06c321c41881f119667fe7c5))
* **golint:** check err on close ([56aafb8](https://github.com/lba-soultec/go-via/commit/56aafb8ae2315ac384acc0a2258e5b8cc91a476e))
* **sha:** use sha as fallback ([3fa1a66](https://github.com/lba-soultec/go-via/commit/3fa1a66ad07aa369bcd2241d5be365894b759909))
* update uploadBackgroundImage endpoint to use theme image URL ([3b4842e](https://github.com/lba-soultec/go-via/commit/3b4842e46410fa5c850f6526f5adfc60f8b6f8d0))
* **version:** use head ref ([d78d3f4](https://github.com/lba-soultec/go-via/commit/d78d3f4007199280f8a3e75d23814f0c569407e0))


### Features

* implement theme management with image upload and retrieval functionality ([7cf96e1](https://github.com/lba-soultec/go-via/commit/7cf96e10d9dd08c823b2edb4d7cc8d94e75c021d))
* remove webauth and first implementation of request made by AI ([4cf67ba](https://github.com/lba-soultec/go-via/commit/4cf67ba5b17b4647ea8b9ffdedd4fc06c588decd))

## [1.0.12](https://github.com/lba-soultec/go-via/compare/v1.0.11...v1.0.12) (2025-08-13)


### Bug Fixes

* **lint:** resolve issues detected by golangci-lint ([b5d4ff4](https://github.com/lba-soultec/go-via/commit/b5d4ff4d0a5a28974a06c32131dbed163e2471a3))

## [1.0.11](https://github.com/lba-soultec/go-via/compare/v1.0.10...v1.0.11) (2025-08-13)


### Bug Fixes

* **migration:** remove existing indexes on old db ([09c5d15](https://github.com/lba-soultec/go-via/commit/09c5d15ccc9f63dfa9ac2050d3adc9b2dad0bb14))

## [1.0.10](https://github.com/lba-soultec/go-via/compare/v1.0.9...v1.0.10) (2025-08-13)


### Bug Fixes

* **creds:** mention the founder ([5e0cfc7](https://github.com/lba-soultec/go-via/commit/5e0cfc7eebe0ceef160f9f064a888d3e7ac8e8b0))

## [1.0.9](https://github.com/lba-soultec/go-via/compare/v1.0.8...v1.0.9) (2025-08-13)


### Bug Fixes

* **naming:** rename to soda ([660f041](https://github.com/lba-soultec/go-via/commit/660f0415d28b0a1c742459f3cb8d3fc74fb2027d))

## [1.0.8](https://github.com/lba-soultec/go-via/compare/v1.0.7...v1.0.8) (2025-08-13)


### Bug Fixes

* **lint:** resolve all issues found by golangci-lint ([c921bc5](https://github.com/lba-soultec/go-via/commit/c921bc5dd7d50c8eaebb3facbe5c80f8e2d140ab))

## [1.0.7](https://github.com/lba-soultec/go-via/compare/v1.0.6...v1.0.7) (2025-08-12)


### Bug Fixes

* **vendoring:** update sum and mod ([3bfb786](https://github.com/lba-soultec/go-via/commit/3bfb786cfca569ba34ce04a4f9a0f5085fb9461d))

## [1.0.6](https://github.com/lba-soultec/go-via/compare/v1.0.5...v1.0.6) (2025-08-12)


### Bug Fixes

* **images:** update image registry to ghcr.io ([6423c01](https://github.com/lba-soultec/go-via/commit/6423c018c79b3dd53d2e60f793742d9c6a6bf71e))
* **image:** use public image ([01abe3a](https://github.com/lba-soultec/go-via/commit/01abe3a775f7657d43c23890ecff50a19528ce73))

## [1.0.5](https://github.com/lba-soultec/go-via/compare/v1.0.4...v1.0.5) (2025-08-12)


### Bug Fixes

* **ns:** add ns ([62427e3](https://github.com/lba-soultec/go-via/commit/62427e3d9d6d618cfefaa0cc023ab5769bc42918))

## [1.0.4](https://github.com/lba-soultec/go-via/compare/v1.0.3...v1.0.4) (2025-08-12)


### Bug Fixes

* **cicd:** add true bool ([151d7d9](https://github.com/lba-soultec/go-via/commit/151d7d9098a09b777802510dc5f1d5a8473fd7a3))

## [1.0.3](https://github.com/lba-soultec/go-via/compare/v1.0.2...v1.0.3) (2025-08-12)


### Bug Fixes

* **deploy:** use kustomize ([e505250](https://github.com/lba-soultec/go-via/commit/e50525012afddacb78a68792ef4c79ce23f53330))

## [1.0.2](https://github.com/lba-soultec/go-via/compare/v1.0.1...v1.0.2) (2025-08-12)


### Bug Fixes

* **cicd:** all branches ([6ca8303](https://github.com/lba-soultec/go-via/commit/6ca83033e9e7dd973c7999b6a7d6a7d4d0314fc5))
* **release:** tag the commit and proceed ([5d0fceb](https://github.com/lba-soultec/go-via/commit/5d0fcebb406997f187f61aaaa871f0dabeb70d66))
* **trigger:** add trigger for build on new runners ([4fe9e68](https://github.com/lba-soultec/go-via/commit/4fe9e68d620cc2bc0734ba06d5ba59d0c5ed6fb8))

## [1.0.2](https://github.com/lba-soultec/go-via/compare/v1.0.1...v1.0.2) (2025-08-12)


### Bug Fixes

* **cicd:** all branches ([6ca8303](https://github.com/lba-soultec/go-via/commit/6ca83033e9e7dd973c7999b6a7d6a7d4d0314fc5))
* **release:** tag the commit and proceed ([5d0fceb](https://github.com/lba-soultec/go-via/commit/5d0fcebb406997f187f61aaaa871f0dabeb70d66))

## [1.0.1](https://github.com/lba-soultec/go-via/compare/v1.0.0...v1.0.1) (2025-08-12)


### Bug Fixes

* add releases ([5c751b4](https://github.com/lba-soultec/go-via/commit/5c751b4ecbc6f216fe313ff7eac1f97ddf3adce9))
* **on:** add create statement ([783811b](https://github.com/lba-soultec/go-via/commit/783811b6beb7e461d1cb11f2c09cadd5e3d02aa6))
* **trigger:** allow trigger ([413a2fb](https://github.com/lba-soultec/go-via/commit/413a2fbe50633f05f166db04dc1d151a57b7b0aa))
* update conditions for build-go-cli, release, and deploy-staging jobs ([f94368f](https://github.com/lba-soultec/go-via/commit/f94368fb1119815c99fafb2aa91f01f32b7eee1a))
* update tag conditions and remove unnecessary job triggers in CI/CD workflow ([68d4778](https://github.com/lba-soultec/go-via/commit/68d4778565aed3a06221efe90effa53e015a99b8))

## [1.0.3](https://github.com/lba-soultec/go-via/compare/v1.0.2...v1.0.3) (2025-08-11)


### Bug Fixes

* update conditions for build-go-cli, release, and deploy-staging jobs ([f94368f](https://github.com/lba-soultec/go-via/commit/f94368fb1119815c99fafb2aa91f01f32b7eee1a))

## [1.0.2](https://github.com/lba-soultec/go-via/compare/v1.0.1...v1.0.2) (2025-08-11)


### Bug Fixes

* add releases ([5c751b4](https://github.com/lba-soultec/go-via/commit/5c751b4ecbc6f216fe313ff7eac1f97ddf3adce9))
* **on:** add create statement ([783811b](https://github.com/lba-soultec/go-via/commit/783811b6beb7e461d1cb11f2c09cadd5e3d02aa6))

## [1.0.1](https://github.com/lba-soultec/go-via/compare/v1.0.0...v1.0.1) (2025-08-11)


### Bug Fixes

* **trigger:** allow trigger ([413a2fb](https://github.com/lba-soultec/go-via/commit/413a2fbe50633f05f166db04dc1d151a57b7b0aa))

# 1.0.0 (2025-08-11)


### Bug Fixes

* **golangci-lint:** fix tftp and statik to check code in action ([5873056](https://github.com/lba-soultec/go-via/commit/58730567e0821c47fdb77092c33eba5b7be8ee6e))


### Features

* add semantic versioning workflow ([716aaac](https://github.com/lba-soultec/go-via/commit/716aaac2bfabc69153801c25c805cec33c909b83))
