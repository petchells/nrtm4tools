# TODO

An unprioritized list of things I'd like to get round to at some point.

- App

  - Use TOML file for configuring notification / source / repo
  - `next_signing_key` verification
  - `--dry-run` option for connect/update to report on what would be done
    -- useful when repo is too far behind server to catch up
  - `validate` command
    - Split into ones that operate remotely-only, and ones that validate
      against a database repository.
    - Print remote server file status
    - Compare remote snapshot with our repo of the same version
    - Consistency check for remote deltas against our historic state
  - Rename snapshot files when hash fails, as with deltas now
  - Support publication of historic states so that mirrors that have lost
    sync with their current server can catch up.
  - Check that the delta file hashes in a notification match previously
    seen values. See [Spec doc §4.3: The mirror client MUST verify... the hashes of each... File](https://htmlpreview.github.io/?https://github.com/mxsasha/nrtmv4/blob/main/draft-ietf-grow-nrtm-v4.html#name-processing-delta-files)
  - Check timestamp of notification file to see if it's stale: §4.4

- Web

  - React f/e to do same things as CLI
  - Add accounts: admin/user for query/update/create privs

- Queries

  - Cross-source
  - History
