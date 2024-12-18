# TODO

An unprioritized list of things I'd like to get round to at some point.

- App

  - Use TOML file for configuring notification / source / repo
  - `next_signing_key` verification
  - `validate` command
  - `validate` command
    - Split into ones that operate remotely-only, and ones that validate
      against a database repository.
    - Print remote server file status
    - Compare remote snapshot with our repo of the same version
    - Consistency check for remote deltas against our historic state
  - `list` command shows details about sources
  - Support the option to overwrite snapshot files when hash fails

- Web

  - React f/e to do same things as CLI
  - Add accounts: admin/user for query/update/create privs

- Queries

  - Cross-source
  - History

- API
  - Support queries by web clients
