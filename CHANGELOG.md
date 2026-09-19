# Changelog

## 0.1.1

- `truenas_file`: creating a new file no longer fails with "Cannot stat file". TrueNAS accepts an
  upload before it writes the file, and the provider now waits for the write to finish before
  reading the file back. A write that fails after the upload is accepted is now reported.
- `truenas_file`: documents that an unset `mode` leaves a new file at `0700`.

## 0.1.0

First release: 35 resources with data sources and list resources, nine actions, and the
`size_bytes` function.
