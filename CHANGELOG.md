# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## [v0.17.0](https://github.com/stefanlogue/meteor/releases/tag/v0.17.0) - 2024-02-06
### Fixed
- Replace `clipboard` package causing panic in published tool

## [v0.16.0](https://github.com/stefanlogue/meteor/releases/tag/v0.16.0) - 2024-02-01
### Added
- custom keymap
- write command to clipboard when aborted or failed

## [v0.15.0](https://github.com/stefanlogue/meteor/releases/tag/v0.15.0) - 2024-01-21
### Added
- prints the commit message on failure for easier retries

## [v0.14.0](https://github.com/stefanlogue/meteor/releases/tag/v0.14.0) - 2024-01-21
### Added
- config option for showing the intro screen

## [v0.13.1](https://github.com/stefanlogue/meteor/releases/tag/v0.13.1) - 2024-01-14
### Changed
- Use `vhs` for README gifs

## [v0.13.0](https://github.com/stefanlogue/meteor/releases/tag/v0.13.0) - 2024-01-14
### Fixed
- shouldn't ask for ticket number when board is `NONE`

## [v0.12.0](https://github.com/stefanlogue/meteor/releases/tag/v0.12.0) - 2024-01-14

## [v0.11.0](https://github.com/stefanlogue/meteor/releases/tag/v0.11.0) - 2024-01-14
### Changed
- migrated `bubbletea` to `huh`

## [v0.10.0](https://github.com/stefanlogue/meteor/releases/tag/v0.10.0) - 2023-11-23
### Changed
- Enabled linkedin announcements

## [v0.9.1](https://github.com/stefanlogue/meteor/releases/tag/v0.9.1) - 2023-11-23
### Fixed
- use correct version of goreleaser (nightly)

## [v0.9.0](https://github.com/stefanlogue/meteor/releases/tag/v0.9.0) - 2023-11-23

## [v0.8.2](https://github.com/stefanlogue/meteor/releases/tag/v0.8.2) - 2023-11-13

## [v0.8.1](https://github.com/stefanlogue/meteor/releases/tag/v0.8.1) - 2023-11-11

## [v0.8.0](https://github.com/stefanlogue/meteor/releases/tag/v0.8.0) - 2023-11-11

## [v0.7.0](https://github.com/stefanlogue/meteor/releases/tag/v0.7.0) - 2023-11-11
### Changed
- Support for POSIX/GNU-style flags, `--version` and `-v` now print the version number

## [v0.6.0](https://github.com/stefanlogue/meteor/releases/tag/v0.6.0) - 2023-11-10
### Added
- Help keys for Coauthor list and message inputs

## [v0.5.0](https://github.com/stefanlogue/meteor/releases/tag/v0.5.0) - 2023-11-09
### Fixed
- Hangs when no config file found

## [v0.4.0](https://github.com/stefanlogue/meteor/releases/tag/v0.4.0) - 2023-11-08

## [v0.3.1](https://github.com/stefanlogue/meteor/releases/tag/v0.3.1) - 2023-11-07
### Changed
- Added `ldflags` to `.goreleaser.yml`
- Moved `version` declaration to outside of `main()`

## [v0.3.0](https://github.com/stefanlogue/meteor/releases/tag/v0.3.0) - 2023-11-06
### Added
- `-v` flag prints installed version
- supports breaking changes

## [v0.2.1](https://github.com/stefanlogue/meteor/releases/tag/v0.2.1) - 2023-11-06

## [v0.2.0](https://github.com/stefanlogue/meteor/releases/tag/v0.2.0) - 2023-11-06
### Added
- Installable via `brew`

## [v0.1.2](https://github.com/stefanlogue/meteor/releases/tag/v0.1.2) - 2023-11-06
### Fixed
- Bug #1: can now deselect coauthors

## [v0.1.1](https://github.com/stefanlogue/meteor/releases/tag/v0.1.1) - 2023-11-06
### Added
- Ability to pick no coauthors

## [v0.1.0](https://github.com/stefanlogue/meteor/releases/tag/v0.1.0) - 2023-11-05
### Added
- Can find config file anywhere up the directory tree

## [v0.0.2](https://github.com/stefanlogue/meteor/releases/tag/v0.0.2) - 2023-11-05
### Added
- Initial release
- Automated releases

### Changed
- Updated flag on `goreleaser`
