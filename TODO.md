# TODO

A loosely prioritized list of things I'd like to get round to at some point.

- App

  - Rename snapshot files when hash fails, as with deltas now
  - Check that the delta file hashes in a notification match previously seen values. See
    [Spec doc §4.3: The mirror client MUST verify... the hashes of each... File](https://htmlpreview.github.io/?https://github.com/mxsasha/nrtmv4/blob/main/draft-ietf-grow-nrtm-v4.html#name-processing-delta-files)
  - Check timestamp of notification file to see if it's stale: §4.4
  - User messages for feedback to CLI and web users. Two types, Response types and Log types:
    - Responses to an action: Error, Warn,...
    - Stream (or sth close to it) log messages to f/e
    - cli should have `-q` option which outputs a one line stdout/err message and an exit code
  - `--dry-run` option for connect/update to report on what would be done
    -- useful to test when repo is too far behind server to catch up, without doing so.
  - `validate` command
    - Split into ones that operate remotely-only, and ones that validate
      against a database repository.
    - Validate all deltas refs in notification, even ones we don't need
    - Compare remote snapshot with our repo of the same version
    - Consistency check for remote deltas against our historic state
    - `next_signing_key` (...when RIPE server publishes one)
  - Use TOML file for configuring notification / source / repo
  - Support publication of historic states so that mirrors that have lost
    sync with their current server can catch up.

- Queries

  - Cross-source
  - History
  - Aggregate functions, reports
  - Data export to Kibana et al? What might be useful formats?

- Web

  - React f/e to do same things as CLI
  - Add accounts: admin/user for query/update/create privs
    - Use OAuth2 or diy? Or both, and make it configurable (depend on TOML step above)
  - New executable to query an API over http
    - Build on privilege component (depend on Add accounts)
    - Need transport layer security as well -- client certs? IP filters?
    - Low prio -- need a few working parts before all this can work.

# Objective

Build a suite of tools for doing :

- `nrtm4`<br>
  Command line tool maintaining a local repo and doing "all the things", as they say.
- `nrtm4serve`<br>
  Web server provides back end services for doing the same things as `nrtm4` in addition,
  eventually, to implementing an API for `nrtm4sclient`
- `nrtm4sclient`<br>
  A CLI executable to communicate with an API provided by `nrtm4serve`
- `nrtm4validator`<br>
  A CLI executable to do repo-less commands on remote servers: A subset of `nrtm4` CLI commands
- `nrtm4mirror`<br>
  A service which complies with the RFC. Not sure about including this
