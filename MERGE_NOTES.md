# Merge notes: jasonlaguidice/master into lwittwer/beta-latest (2026-11-15)

Branch `sync-jason-20261115` on `mrgnw/imessage`, based on `lwittwer/beta-latest` @ `85d57b37`.
Jason side: `jason/master` @ `595ce7e9`. Merge base `905003c2` (2026-05-20). lwittwer was 543 ahead, 40 behind.

## Method

A trial `git merge jason/master` produced 61 conflicting files, including a 5,600-line divergence in
`pkg/rustpushgo/src/lib.rs`. The cause: both forks pulled the same upstream (lrhodin/corten-matrix, now
private) but in different shapes. lwittwer merged upstream directly through 2026-10-09. Jason re-ported the
same upstream commits by hand into a multi-user, module-split layout (`aps_watchdog.rs`, `panic_guard.rs`,
`facetime_notify.rs`, per-login StatusKit keys, settings in the user DB).

Every corten-matrix commit hash referenced in jason's tree (19 of them, for example `cd643d57`, `0f8623b4`,
`993ec4e3`, `af1fe3e1`, `a4b512ae`) is already an ancestor of lwittwer/beta-latest. A union merge would have
meant hand-reconciling two implementations of the same upstream code. So jason's own work was cherry-picked
commit by commit instead (`-x`, original author and date kept), and the upstream re-ports were not taken.
Because of this, `git log beta-latest..jason/master` will still list all 40 commits. A future sync should
start from this table, not from the commit count.

## The 40 jason-only commits

Disposition key:
- TAKEN: cherry-picked.
- PORTED: re-applied by hand onto lwittwer code, jason kept as author.
- IN-LW: lwittwer already has it, either the same fix from upstream or an equivalent.
- SKIP: not taken, reason given.
- DEFER: worth porting, not done in this pass.

| Commit | Subject | Disposition |
|-|-|-|
| fd462102 | Revise for multi-user & add docker build flow | SKIP: multi-user-in-one-process refactor; lwittwer runs one process per account (`bridge-all`). Docker/compose are Linux-only and depend on jason's multi-user changes |
| ecf838d6 | Re-add ca-certs to docker container | SKIP (Docker) |
| c955b9ef | Add /data directory to working path | SKIP (Docker) |
| a8713b33 | Move settings from config to user database | SKIP: replaces lwittwer's config model and rewrites the install scripts (about 1,000 lines). Open decision |
| f0c9dc23 | Undo sudo removal from dev environment | SKIP: depends on a8713b33's script rewrite |
| 42476c2f | Add manual dispatch | SKIP (jason's Docker CI) |
| daa6f0f4 | Move config out of entrypoint | SKIP (Docker) |
| cbc42dfe | Remove mautrix branding (sorry tulir) | SKIP: branding (`cmd/matrix-rustpush`); lwittwer uses `cmd/corten-matrix` |
| b70a5936 | Add images to latest tag | SKIP (jason's Docker CI) |
| a5f4a13f | Add video transcoding and heic conversion commands | SKIP: needs a8713b33 (per-user settings) |
| 772f22b6 | Bypass permission check when running as root | TAKEN `a9802ac6` |
| 990d630a | Add logging to media conversion steps | TAKEN `23048677`; conflict resolved, see below |
| c96e51a0 | Remove more singletons from rustpush | SKIP: per-login state filenames for multi-user; would orphan existing `statuskit-state.plist` / `facetime-state.plist` |
| 7aabd86e | Add catch for crash when CloudKit profile cannot be retrieved | IN-LW: lwittwer wraps the same `get_record` in a panic catch (`lib.rs` "inline ShareProfile fetch panicked") |
| 822ca32b | Add open_absinthe=info to default RUST_LOG | PORTED `d3ffa9f5` into lwittwer's `DEFAULT_RUST_LOG` |
| 6bea02f8 | Include x86_64 binary blobs for ARM64 environments | SKIP: Linux arm64 open-absinthe emulation; lwittwer deleted `open-absinthe/build.rs` (native NAC on macOS only, `8c2d9ea9`) |
| 73bdd924 | Fix arm64 build cross-compile issue | SKIP (same reason) |
| e1d4e56c | Add gap between .data and .text for arm64 | SKIP (same reason) |
| db19560d | Change SQLite placeholders for standard Postgres type | IN-LW (patch reverse-applies cleanly) |
| e1ab8fe9 | Disable release workflow (for now) | SKIP (jason's CI) |
| e9e1f407 | Bump readme | SKIP (jason's README) |
| ef214536 | Add correct libheif to docker image | SKIP (Docker) |
| 0ea1ee0c | Fix PostgreSQL compatibility in cloud backfill and shared profile | IN-LW: `sqlBlobType` plus dialect-aware column migrations in `cloud_backfill_store.go` |
| 564c3516 | Bump go version for vulnerability fix | IN-LW (`git cherry` patch-id match) |
| 5a7d196d | Adjust build tagging for branches | SKIP (jason's Docker CI) |
| 5d083612 | Port four new config knobs from upstream/master | IN-LW: `max_attachment_size_mb`, `url_previews_in_backfill`, `read_receipts`, `typing_notifications` are all in `config.go` |
| 1b82804e | Port contacts reliability and group roster fixes | IN-LW (upstream `993ec4e3`, `108a196e`, …) |
| cd17dbcd | Port FaceTime bug fixes from upstream/master | IN-LW (upstream `a4b512ae` notify-only FaceTime and later) |
| 6957b747 | June changes roll-up from lrhodin/corten | IN-LW (rlimit/memlimit, APNs watchdog, wedge recovery are all upstream commits lwittwer merged) |
| 3d90823b | restore identity_for_reregister capture | SKIP: fixes jason's own refactor; lwittwer has no such variable |
| ddeba1b5 | README simplification | SKIP (jason's README) |
| 36358aff | Add v prefix to container image tag | SKIP (Docker CI) |
| 47068c9e | Remove extra version tag as well | SKIP (Docker CI) |
| 5f85500c | StatusKit notifications in DM room topics (default) | DEFER: portable, but defaults to `topic` and so changes lwittwer's notice behaviour. Open decision |
| ab4c015d | Release workflow fix for rebrand | SKIP (jason's CI) |
| 8d5cf801 | Replace network icon w/ Rustpush/OpenBubbles icon | SKIP: branding; lwittwer chose its icon in `451996cc` |
| 3b12dfb5 | Add resiliency to CardDAV sync failures | TAKEN `cd92f45e`; conflict resolved, see below |
| d386645d | Add identity check and switch commands (#3) | PORTED `a8ead1a6`: `list-identities`, `set-identity` |
| 3f4c1fe7 | Remove custom logout path, revert to bridgev2 logout | SKIP: kept lwittwer's `logout` (numbered menu, scrubs session backup). Jason's per-account scrub targets multi-user; see credential hardening item 2 |
| 595ce7e9 | Sync/upstream 2026 09 (#5) | IN-LW for the upstream content. Jason-only extras: escrow-record fallback (IN-LW as the Makefile `passcode_generation` serde-default patch), `keychain-join` command (DEFER), `backfill_cutoff.go` (SKIP: only needed on mautrix-go ≥ v0.30; lwittwer is on v0.29), mautrix-go v0.31 bump (SKIP), multi-user scoping (SKIP) |

## Conflicts and decisions

- `990d630a` (`client.go`): lwittwer routes every video transcode through `transcodeToMP4` (shared semaphore
  and retry). Jason added log lines at three inline ffmpeg call sites that lwittwer no longer has. Decision:
  the two log lines ("Starting video transcoding", "Video remux failed, falling back to re-encode") go inside
  `transcodeToMP4`, so every path logs them. The `heic.go` line applied cleanly.
- `3b12dfb5` (`external_carddav.go`): kept lwittwer's sanitised log line (`sanitizeURLError`, `logSafeURL`)
  and added jason's `fetchFailures` counter and the all-books-failed early return.
- `d386645d`: rewritten onto lwittwer. Every `c.handle` and `c.allHandles` read in non-test
  `pkg/connector` files (71 changed lines) now calls `getHandle()` / `getAllHandles()`. `Connect` sets both
  through `setIdentity`. The commands sit in a new `pkg/connector/identity_commands.go` carrying the AGPL-3.0
  header from jason's `commands.go`, so they are not presented as MPL. `parseHandleSelection` and
  `ForceReregisterIdentity` already existed in lwittwer.
- `822ca32b`: lwittwer moved the default into the `DEFAULT_RUST_LOG` constant. `open_absinthe=info` was
  appended there.

## Licence findings

Files at the root: only `LICENSE`, which is MPL-2.0. The upstream maintainer (mackid1993) replaced AGPL-3.0
with MPL-2.0 in `f84bc7eb` and `ee81bf19` on 2026-06-26. The licence text was not changed in this merge.

Per-file headers in the merged tree (tracked source, not counting `third_party/`, which is cloned at build
time):

- 30 files carry the MPL-2.0 header ("This Source Code Form is subject to the terms of the Mozilla Public
  License, v. 2.0"). Examples: `pkg/connector/client.go`, `commands.go`, `login.go`, `identity_store.go`,
  `cmd/corten-matrix/login_cli.go`.
- 22 files carry the AGPL-3.0 header:
  - 21 restored by upstream `875e8eaf` ("keep Tulir Asokan's files under AGPL-3.0"): `cmd/corten-matrix/main.go`,
    `imessage/interface.go`, `imessage/struct.go`, `imessage/tapback.go`, `ipc/ipc.go`, and all
    `imessage/mac/*.go|.h|.m`.
  - 1 new from this merge: `pkg/connector/identity_commands.go` (jason's code).
- 144 files carry no licence header.

Does merging AGPL code into the MPL tree create an obligation?

- Yes. Jason's fork is AGPL-3.0: his `LICENSE`, his file headers, and he never relicensed. His contributions
  stay AGPL-3.0. That covers `identity_commands.go` plus the small hunks now inside MPL-headed files:
  - the accessors in `client.go`
  - the log lines in `client.go` and `heic.go`
  - 11 lines in `external_carddav.go`
  - 3 lines in `imessage/struct_test.go`
  - the RUST_LOG string in `lib.rs`
- MPL-2.0 §3.3 lets MPL files be combined with GPL-family code in a "Larger Work". It does not relicense the
  AGPL hunks.
- In practice nothing changes for the binary. `cmd/corten-matrix/main.go` and the `imessage/` packages were
  already AGPL-3.0, so any distributed or network-served `corten-matrix` binary already falls under AGPL-3.0
  terms. AGPL §13 requires offering source to users who interact with a modified version over a network. A
  public fork satisfies that.
- Unverified: at the merge base every connector file was AGPL-3.0 with many authors (Tulir Asokan 934
  commits, Ludvig Rhodin 566, David 502, Tae Hagen 268, …). Whether upstream's later MPL relicense of those
  files had every contributor's consent cannot be checked from the repo. Treat the tree as AGPL-3.0 overall
  if in doubt.
- Open decision: either leave the in-file jason hunks under MPL-headed files, as now, or add a notice. No
  header was changed by guesswork.

## Build and test (macOS 27.2, Apple Silicon, Go 1.27.2, Rust 1.97.1)

- `make rust`: PASS, after two fixes to the environment (not the code):
  - `protoc` was missing. `brew install protobuf` is what the Makefile's `check-deps` would run.
  - Proc-macro dylibs fail to load ("mis-aligned LINKEDIT string pool") when cargo's default release
    `strip=debuginfo` runs. This happened with and without sccache, and on both the internal and external
    disks. It is a toolchain/OS issue, not caused by the merge. Workaround used:
    `RUSTC_WRAPPER= CARGO_PROFILE_RELEASE_STRIP=none CARGO_BUILD_BUILD_DIR=… make rust`.
- `go build ./...` and `go build ./cmd/corten-matrix`: PASS, with the Makefile's `CGO_CFLAGS` and
  `CGO_LDFLAGS`. Codesign and install were not run.
- `go vet ./...`: exit 1, only 11 "possible misuse of unsafe.Pointer" in the generated uniffi file
  `pkg/rustpushgo/rustpushgo.go`, which this merge did not touch. Nothing in hand-written code.
- `go test ./...`: PASS (8 packages, run with `XDG_DATA_HOME` and `XDG_CONFIG_HOME` pointed at a scratch dir).
- `cargo test --release --no-default-features --features nac-apple-framework` (pkg/rustpushgo): PASS,
  12 unit tests and 2 `password_state_restore` tests.
- Not run: the bridge itself, login, setup, `make install`. The new commands have no tests and were not
  exercised.

GitHub Actions on the fork: the API reports Actions as enabled. Pushes start nothing, though:
- `build-release-binaries.yml` is `workflow_dispatch` only and needs `OPENCIDER_READ_DEPLOY_KEY` (not set).
- `codex-review-on-ready.yml` runs on `pull_request_target` only.

## Credential audit (read-only)

### "preserve legacy password sharing state" (`51111de9`, Lucas Wittwer, 2026-10-09)

- "Password sharing state" is rustpush's `passwords::PasswordState`: the Apple Passwords-app shared password
  groups and pending invitations. It has nothing to do with the Apple ID password. It is persisted to
  `$XDG_DATA_HOME/corten-matrix/passwords-state.plist`.
- "Legacy" means a plist written by older rustpush. That version had `token_registered: bool` where the
  current one has `my_token_registered: Option<[u8;32]>`. The commit adds a Makefile `rp_patch` that marks
  the new field `#[serde(default)]`, so an old plist still loads (groups and invites kept, APNs token
  re-registered). It also adds `pkg/rustpushgo/tests/password_state_restore.rs`.
- In this tree the Passwords client is dormant. `get_passwords_client` is exported over FFI, but no Go code
  calls it (`rg PasswordsClient pkg/connector` finds nothing). So the plist should never be created.
  Unverified at runtime.

### Apple ID password: what is stored, and where

Both branches store the same fields:

| Location | Fields | Format |
|-|-|-|
| Bridge DB `user_login.metadata` (`UserLoginMetadata`, `pkg/connector/dbmeta.go`) | `account_hashed_password_hex`, `account_persist_blob` (binary plist with `PersistAccountData.hashed_password` and GSA tokens), `account_pet`, `account_spd_base64`, `mme_delegate_json`, `ids_identity` (device private keys), `aps_state`, `ids_users` | plaintext JSON in SQLite |
| `session.json` (`PersistedSessionState`, `pkg/connector/identity_store.go`) | the same set, as a backup that survives DB resets | plaintext JSON |

- `hashed_password` is SHA-256 of the raw password with no salt
  (`icloud-auth/src/client.rs:818`). It is what the client feeds into SRP/PBKDF2, so it works as a password
  equivalent for GSA logins. Unsalted SHA-256 can also be dictionary-attacked offline.
- None of it is encrypted. Only the CardDAV app password is encrypted (AES-256-GCM, key in a 0600 key file,
  `carddav_crypto.go`).
- lwittwer vs jason:
  - lwittwer: `~/.local/share/corten-matrix/session.json`, written atomically (temp file, 0600, rename),
    directory 0700.
  - jason: `~/.local/share/mautrix-imessage/session.json`, written with `os.WriteFile(…, 0600)`, which leaves
    the mode of an already-existing file as it was. In jason's multi-user mode it is one bridge-wide file,
    last writer wins.
  - The DB fields are identical in both.
- Bridge DB file mode: unverified. SQLite creates it under the process umask, usually 0644, and nothing in
  `pkg/cli` chmods it.

### Is it reused or shared?

- Reused for refresh, yes:
  - `Client::refresh_pet_token` (`lib.rs` around line 2160) re-runs `login_email_pass(username, hashed_password)`
    every 12 h.
  - Upstream `TokenProvider::get_token` (`icloud-auth/client.rs:461`) does the same automatically.
  - Without the stored hash, the PET expires and CardDAV and CloudKit need a manual re-login. Whether IDS
    messaging keeps working without it is unverified.
- `bridge-all` does not share it. The second account runs as a separate process with its own
  `XDG_DATA_HOME` (`…/corten-matrix-1`), its own session.json and its own DB (`pkg/cli/cli.go:512`). Both run
  as the same OS user, so each can read the other's files.
- No other process or device gets it. SRP sends proofs, not the hash. On macOS NAC/anisette is native
  (`nac-apple-framework`). The remote anisette host `ani.sidestore.io` is compiled only with the Linux
  `anisette-remote-v3` feature.

### Privacy scrubber

- `logSafeHandle` / `logSafeHandles` / `logSafeURL` hash handles and strip URLs in connector logs
  (`client.go:2134`). It is off when `debug_disable_privacy` is set.
- No log line in `pkg/connector`, `cmd`, `pkg/cli` or `lib.rs` prints the password, hash, PET or delegate.
  They log booleans only, for example "PET token available: true".
- The scrubber does not cover the DB or session.json; credentials sit there in clear.
- rustpush's own debug/trace logging under `RUST_LOG=debug` was not audited (unverified).

### `corten-matrix update` (self-update)

Not in this source tree; there is no `update` subcommand in `cmd/corten-matrix/main.go` or `pkg/cli/cli.go`.
README commit `154ac6a9` says it ships "only on the official prebuilt binaries". Those are built by
`build-release-binaries.yml` from the private `mackid1993/OpenCider` repo (deploy key), so their contents
cannot be audited. To disable it in our fork, build from source and never run a release binary. Nothing in
the source needs removing. Releases are fetched from GitHub
(`mackid1993/corten-matrix-releases`, `lrhodin/corten-matrix/releases`).

### Outbound traffic other than Apple

- No telemetry SDKs (no sentry, posthog or analytics).
- Non-Apple traffic:
  - your Matrix homeserver
  - CardDAV providers you configure
  - URL previews: the bridge fetches links found in messages, which exposes your IP to those sites;
    `url_previews_in_backfill` gates only the backfill case
  - Linux only: `ani.sidestore.io`

### Recommended hardening (not implemented), ranked by value against effort

1. Build from source only and never run `corten-matrix update` or a release binary. This removes the
   self-updater and the unauditable OpenCider code. Effort: none.
2. Scrub session.json on every logout path. `IMClient.LogoutRemote` (`client.go:2049`) only calls
   `Disconnect`, so a provisioning-API or bridgev2 logout leaves the hash, PET and IDS keys on disk. Only the
   chat `logout` command scrubs. Fix: call `clearLocalSessionBackup` (`commands.go:526`) from `LogoutRemote`
   when the logging-out account owns the file (jason's `clearSessionBackupForAccount` shows the DSID check).
   Effort: small.
3. Tighten file modes at startup: chmod the bridge DB and its `-wal`/`-shm` files to 0600 and the data dir to
   0700. Place: `IMConnector.Start` or `cmd/corten-matrix/ensure_config.go`. Effort: small.
4. Keep one copy of the hash, not two. `account_hashed_password_hex` duplicates what `account_persist_blob`
   already holds. Make the hex field a fallback that is cleared once the blob round-trips
   (`login.go:796`, `client.go:1188`, `identity_store.go`). Effort: small to medium.
5. Encrypt credentials at rest with a key held in the macOS Keychain. Cover the `UserLoginMetadata` and
   session.json fields listed above. Reuse the AES-GCM code in `carddav_crypto.go`, but load the key from the
   Keychain instead of a key file. Check that a launchd agent can read the login keychain while the screen is
   locked (unverified). Effort: medium.
6. Optional "no stored password" mode. Do not persist `hashed_password`, and accept a re-login when the PET
   expires (about 24 h; the exact lifetime is unverified). Touch `persist_data_from_spd`,
   `refresh_pet_token`, `login.go:796`. Effort: medium. Cost: CardDAV and CloudKit stop working until you log
   in again.
7. Add a config switch for live URL previews (`urlpreview.go`, `fetchURLPreview`) if leaking your IP to
   linked sites matters. Effort: small.

## Open decisions for the user

- Multi-user and per-user settings (a8713b33, a5f4a13f, c96e51a0) and the Docker flow: take them, or stay on
  lwittwer's one-process-per-account model?
- StatusKit topic style (5f85500c): port it with lwittwer's notice behaviour as the default?
- `keychain-join` (from 595ce7e9): port it? It asks for a device passcode in chat.
- Licence labelling of jason's hunks inside MPL-headed files.
- Which credential-hardening items to implement.
