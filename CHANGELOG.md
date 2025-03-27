# Changelog

## [0.7.0](https://github.com/elct9620/terraform-provider-lambdalabs/compare/v6.0.2...v0.7.0) (2025-03-27)


### Features

* Add computed attributes for filesystem resource based on CreateFileSystem API ([39082f5](https://github.com/elct9620/terraform-provider-lambdalabs/commit/39082f518dca619a8aced8b37562fd3d24b57fa8))
* Add create filesystem API implementation with comprehensive tests ([e87c8c4](https://github.com/elct9620/terraform-provider-lambdalabs/commit/e87c8c4fcfade3b1d58ec71ea0ca0df2dfa13779))
* Add delete filesystem API with comprehensive test cases ([bbb2a71](https://github.com/elct9620/terraform-provider-lambdalabs/commit/bbb2a71e32e9e91b23d4c61bec09840bfad036aa))
* Add filesystem data provider with initial implementation and tests ([950dd87](https://github.com/elct9620/terraform-provider-lambdalabs/commit/950dd874de071ae0a922a3bd25a674d124e4af73))
* Add filesystem package for Lambda Labs storage operations ([f0a3c50](https://github.com/elct9620/terraform-provider-lambdalabs/commit/f0a3c50e9b604239d39d4f75977fe4167f94cd24))
* Add filesystem resource and corresponding test for provider ([97e4466](https://github.com/elct9620/terraform-provider-lambdalabs/commit/97e44663ef1bf968842f429ec22f2beac654783a))
* Add filter support for images data source with region, family, and architecture filters ([88d1644](https://github.com/elct9620/terraform-provider-lambdalabs/commit/88d16447df34e0a8d940b939c9ec26a4a62b60cd))
* Add firewall data source and test files for provider ([8797a1e](https://github.com/elct9620/terraform-provider-lambdalabs/commit/8797a1ecfc54058794c7b2fd02f2e7116da04272))
* Add firewall data source implementation and tests ([90ea13e](https://github.com/elct9620/terraform-provider-lambdalabs/commit/90ea13e7d76c7628d9c507056859b75b8504f77a))
* Add firewall management package for LambdaLabs infrastructure ([c30a7c5](https://github.com/elct9620/terraform-provider-lambdalabs/commit/c30a7c54e88f1e692de5fa68d20ecd91c1048dc2))
* Add image data provider with initial implementation and tests ([0e24f77](https://github.com/elct9620/terraform-provider-lambdalabs/commit/0e24f7729a05bba2719e149fcf24f31a621a1561))
* Add images data source with tests for LambdaLabs provider ([6c39b61](https://github.com/elct9620/terraform-provider-lambdalabs/commit/6c39b6139173a7be9207d4017401517820597996))
* Add instance types data provider for internal configuration ([8941d5a](https://github.com/elct9620/terraform-provider-lambdalabs/commit/8941d5a29319cf66536d9382eae76d031804b3d1))
* Add instance types data source to provider ([b29353a](https://github.com/elct9620/terraform-provider-lambdalabs/commit/b29353a6e38bdf5aa488d06f6a7a2610ab5090b0))
* Add LambdaLabs image generation package with initial implementation and tests ([9c19ffe](https://github.com/elct9620/terraform-provider-lambdalabs/commit/9c19ffe1fabbde9a7261b7da8a3347db777b8999))
* Add List Available Images API support ([cdc77bc](https://github.com/elct9620/terraform-provider-lambdalabs/commit/cdc77bc939593875ea8e9be369edd46ca4c643c0))
* Add List Filesystems API implementation with tests ([59fbb01](https://github.com/elct9620/terraform-provider-lambdalabs/commit/59fbb01be01e86adf3524a179c4feb5adde66ccd))
* Add ListFirewallRules API with tests to Lambda Labs client ([0086efc](https://github.com/elct9620/terraform-provider-lambdalabs/commit/0086efcba4f35c95ac73984da897dd79bcf1cff0))
* Add ListInstanceTypes API and related data structures ([a3aac8f](https://github.com/elct9620/terraform-provider-lambdalabs/commit/a3aac8f39e8d4590344ed98ae86c6b9a68118327))
* Add optional region filter for instance types data source ([4fd747a](https://github.com/elct9620/terraform-provider-lambdalabs/commit/4fd747a35cbb99d28b2df7a4e0edb5b3f576e158))
* Add Replace Inbound Firewall Rule method to Lambda Labs client ([dccd02a](https://github.com/elct9620/terraform-provider-lambdalabs/commit/dccd02a949c37d3122f01e2d5108dd0975430825))
* Implement filesystem data source with filtering and tests ([ffee4d9](https://github.com/elct9620/terraform-provider-lambdalabs/commit/ffee4d99ee81c109ae13997d4ea6dda26c185b81))
* Implement filesystem resource for Lambdalabs provider ([f9eaf6d](https://github.com/elct9620/terraform-provider-lambdalabs/commit/f9eaf6d99ddbd6e77adb784f07ec1f13e247d010))
* Implement instance types data source based on ssh_key_data.go ([ed0d3d3](https://github.com/elct9620/terraform-provider-lambdalabs/commit/ed0d3d3f647cbfc9046376018f69e9fa5ed28310))


### Bug Fixes

* Check error when unmarshaling JSON in filesystem test ([33b72a6](https://github.com/elct9620/terraform-provider-lambdalabs/commit/33b72a6c32c2bcd5d6769781ab1275dd8f99fd3c))
* Correct client initialization function names in firewall test ([315c9f8](https://github.com/elct9620/terraform-provider-lambdalabs/commit/315c9f87f40e1f982415817306ea051d98f35b3c))
* Correct field name from `Id` to `ID` in filesystem data source ([1462a4c](https://github.com/elct9620/terraform-provider-lambdalabs/commit/1462a4cc35b1d140acf6beae1cbbd86142dcdbcb))
* Correct filesystem endpoint paths for list, create, and delete operations ([6f56f9b](https://github.com/elct9620/terraform-provider-lambdalabs/commit/6f56f9b7212a24861eecbf0cc6ca12b98c72f02f))
* Handle write error in firewall_test.go by using blank identifiers ([3bf91dd](https://github.com/elct9620/terraform-provider-lambdalabs/commit/3bf91dd8e5b89513c5d911770f0c56f3f7ab7141))
* Remove unused body variable in firewall test ([dec36cf](https://github.com/elct9620/terraform-provider-lambdalabs/commit/dec36cf7c7eba3f7a3fccf281fc0a4b0109471ff))
* Update filesystem API endpoints from `/file-systems` to `/filesystems` ([7d01517](https://github.com/elct9620/terraform-provider-lambdalabs/commit/7d015175e9c569a89907f6a1051f4d26a9bf8b64))
* Update filesystem mock server endpoints in test file ([8e5aabd](https://github.com/elct9620/terraform-provider-lambdalabs/commit/8e5aabd0522543767f75a109ff84a6c342f4b24e))


### Miscellaneous Chores

* release 0.7.0 ([66b7608](https://github.com/elct9620/terraform-provider-lambdalabs/commit/66b760837f332f87372b894e05a16823015438f4))
* release 6.0.1 ([bf2112f](https://github.com/elct9620/terraform-provider-lambdalabs/commit/bf2112f7826d17c8beff72b703b652fd8e691dd8))
* release 6.0.2 ([4ccd7d6](https://github.com/elct9620/terraform-provider-lambdalabs/commit/4ccd7d6a6cc992bfa27a3f1955c657f92149f73d))

## [6.0.2](https://github.com/elct9620/terraform-provider-lambdalabs/compare/v6.0.1...v6.0.2) (2025-03-13)


### Miscellaneous Chores

* release 6.0.2 ([4ccd7d6](https://github.com/elct9620/terraform-provider-lambdalabs/commit/4ccd7d6a6cc992bfa27a3f1955c657f92149f73d))

## [6.0.1](https://github.com/elct9620/terraform-provider-lambdalabs/compare/v0.6.0...v6.0.1) (2025-03-13)


### Miscellaneous Chores

* release 6.0.1 ([bf2112f](https://github.com/elct9620/terraform-provider-lambdalabs/commit/bf2112f7826d17c8beff72b703b652fd8e691dd8))

## [0.6.0](https://github.com/elct9620/terraform-provider-lambdalabs/compare/v0.5.0...v0.6.0) (2025-03-13)


### Features

* Add computed attributes for filesystem resource based on CreateFileSystem API ([39082f5](https://github.com/elct9620/terraform-provider-lambdalabs/commit/39082f518dca619a8aced8b37562fd3d24b57fa8))
* Add create filesystem API implementation with comprehensive tests ([e87c8c4](https://github.com/elct9620/terraform-provider-lambdalabs/commit/e87c8c4fcfade3b1d58ec71ea0ca0df2dfa13779))
* Add delete filesystem API with comprehensive test cases ([bbb2a71](https://github.com/elct9620/terraform-provider-lambdalabs/commit/bbb2a71e32e9e91b23d4c61bec09840bfad036aa))
* Add filesystem data provider with initial implementation and tests ([950dd87](https://github.com/elct9620/terraform-provider-lambdalabs/commit/950dd874de071ae0a922a3bd25a674d124e4af73))
* Add filesystem package for Lambda Labs storage operations ([f0a3c50](https://github.com/elct9620/terraform-provider-lambdalabs/commit/f0a3c50e9b604239d39d4f75977fe4167f94cd24))
* Add filesystem resource and corresponding test for provider ([97e4466](https://github.com/elct9620/terraform-provider-lambdalabs/commit/97e44663ef1bf968842f429ec22f2beac654783a))
* Add filter support for images data source with region, family, and architecture filters ([88d1644](https://github.com/elct9620/terraform-provider-lambdalabs/commit/88d16447df34e0a8d940b939c9ec26a4a62b60cd))
* Add image data provider with initial implementation and tests ([0e24f77](https://github.com/elct9620/terraform-provider-lambdalabs/commit/0e24f7729a05bba2719e149fcf24f31a621a1561))
* Add images data source with tests for LambdaLabs provider ([6c39b61](https://github.com/elct9620/terraform-provider-lambdalabs/commit/6c39b6139173a7be9207d4017401517820597996))
* Add instance types data provider for internal configuration ([8941d5a](https://github.com/elct9620/terraform-provider-lambdalabs/commit/8941d5a29319cf66536d9382eae76d031804b3d1))
* Add instance types data source to provider ([b29353a](https://github.com/elct9620/terraform-provider-lambdalabs/commit/b29353a6e38bdf5aa488d06f6a7a2610ab5090b0))
* Add LambdaLabs image generation package with initial implementation and tests ([9c19ffe](https://github.com/elct9620/terraform-provider-lambdalabs/commit/9c19ffe1fabbde9a7261b7da8a3347db777b8999))
* Add List Available Images API support ([cdc77bc](https://github.com/elct9620/terraform-provider-lambdalabs/commit/cdc77bc939593875ea8e9be369edd46ca4c643c0))
* Add List Filesystems API implementation with tests ([59fbb01](https://github.com/elct9620/terraform-provider-lambdalabs/commit/59fbb01be01e86adf3524a179c4feb5adde66ccd))
* Add ListInstanceTypes API and related data structures ([a3aac8f](https://github.com/elct9620/terraform-provider-lambdalabs/commit/a3aac8f39e8d4590344ed98ae86c6b9a68118327))
* Add optional region filter for instance types data source ([4fd747a](https://github.com/elct9620/terraform-provider-lambdalabs/commit/4fd747a35cbb99d28b2df7a4e0edb5b3f576e158))
* Implement filesystem data source with filtering and tests ([ffee4d9](https://github.com/elct9620/terraform-provider-lambdalabs/commit/ffee4d99ee81c109ae13997d4ea6dda26c185b81))
* Implement filesystem resource for Lambdalabs provider ([f9eaf6d](https://github.com/elct9620/terraform-provider-lambdalabs/commit/f9eaf6d99ddbd6e77adb784f07ec1f13e247d010))
* Implement instance types data source based on ssh_key_data.go ([ed0d3d3](https://github.com/elct9620/terraform-provider-lambdalabs/commit/ed0d3d3f647cbfc9046376018f69e9fa5ed28310))


### Bug Fixes

* Check error when unmarshaling JSON in filesystem test ([33b72a6](https://github.com/elct9620/terraform-provider-lambdalabs/commit/33b72a6c32c2bcd5d6769781ab1275dd8f99fd3c))
* Correct field name from `Id` to `ID` in filesystem data source ([1462a4c](https://github.com/elct9620/terraform-provider-lambdalabs/commit/1462a4cc35b1d140acf6beae1cbbd86142dcdbcb))
* Correct filesystem endpoint paths for list, create, and delete operations ([6f56f9b](https://github.com/elct9620/terraform-provider-lambdalabs/commit/6f56f9b7212a24861eecbf0cc6ca12b98c72f02f))
* Update filesystem API endpoints from `/file-systems` to `/filesystems` ([7d01517](https://github.com/elct9620/terraform-provider-lambdalabs/commit/7d015175e9c569a89907f6a1051f4d26a9bf8b64))
* Update filesystem mock server endpoints in test file ([8e5aabd](https://github.com/elct9620/terraform-provider-lambdalabs/commit/8e5aabd0522543767f75a109ff84a6c342f4b24e))
