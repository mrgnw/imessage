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
| 5f85500c | StatusKit notifications in DM room topics (default) | PORTED `570f3596` on `harden-credentials`: `statuskit_notification_style`, default `title` (lwittwer's moon), `topic` optional. Jason's per-user override, install prompts and `notice` style not taken |
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
  - the StatusKit topic option (`570f3596`, branch `harden-credentials`): `statusKitTopic`, the topic branch
	in `OnStatusUpdate` and the `statuskit_notification_style` key
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
- Branch `harden-credentials`, same environment:
  - `make rust`, `go build ./...` and `go test ./...` (9 packages): PASS.
  - `go vet ./pkg/connector ./cmd/corten-matrix`: clean.
  - `cargo test`: 13 unit tests and 2 `password_state_restore` tests, PASS.

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

### Recommended hardening, ranked by value against effort

Status as of branch `harden-credentials` (based on `sync-jason-20261115`). Each item names its commit.

1. Build from source only and never run `corten-matrix update` or a release binary. This removes the
   self-updater and the unauditable OpenCider code. Effort: none.
   Status: nothing to change in code; it is an operating rule.
2. Scrub session.json on every logout path. `IMClient.LogoutRemote` (`client.go:2049`) only calls
   `Disconnect`, so a provisioning-API or bridgev2 logout leaves the hash, PET and IDS keys on disk. Only the
   chat `logout` command scrubs. Fix: call `clearLocalSessionBackup` (`commands.go:526`) from `LogoutRemote`
   when the logging-out account owns the file (jason's `clearSessionBackupForAccount` shows the DSID check).
   Effort: small.
   Status: done, `16b23a53`. `clearSessionBackupForAccount` (`identity_store.go`) removes session.json,
   identity.plist and trustedpeers.plist when session.json's DSID or IDS users match the login. It runs
   after `Disconnect`, which waits for the final state save. Test: `TestClearSessionBackupForAccount`.
3. Tighten file modes at startup: chmod the bridge DB and its `-wal`/`-shm` files to 0600 and the data dir to
   0700. Place: `IMConnector.Start` or `cmd/corten-matrix/ensure_config.go`. Effort: small.
   Status: done, `b3712b43`, in `cmd/corten-matrix/filemodes.go`. The connector never sees the DB URI, so
   it runs in `main.go`: before `Init` and after `Start` (a DB created on first run), and after `Init` in
   `init-db`. SQLite gives new `-wal`/`-shm` files the DB's mode, so they stay 0600. In-memory DSNs,
   non-SQLite configs and missing files are no-ops. Only `$XDG_DATA_HOME/corten-matrix` is set to 0700,
   not an arbitrary DB parent directory. Tests: `TestTightenFileModes*`.
4. Keep one copy of the hash, not two. `account_hashed_password_hex` duplicates what `account_persist_blob`
   already holds. Make the hex field a fallback that is cleared once the blob round-trips
   (`login.go:796`, `client.go:1188`, `identity_store.go`). Effort: small to medium.
   Status: done behind `network.drop_redundant_password_hash`, default false, `a63ce667`. The key sits under
   `network:` like every other connector option; bridgev2 owns `bridge:`.
   - A small binary-plist reader (`password_hash.go`) pulls `account.hashed_password` out of the blob,
	 only when the blob version is 1 and its ADSID matches the login. These are the same checks
	 `restore_token_provider` applies.
   - With the flag on, the hex field is cleared only when the blob holds the identical hash.
   - Restore, `validateCloudKitRestoreState` and every blob update use the hex when set, else the
	 blob's hash. If a new blob ever lacks the hash, the hex copy is written back first.
   - Refresh (`refresh_pet_token`) reads the hash from the in-memory account, and restore fills that
	 from the blob. So refresh behaves the same either way.
   - Why the flag: the empty-hex restore was proven only by unit tests, not by a live restart or a
	 12 h PET refresh. A Rust test pins `encode_persist_blob` output byte for byte, and the Go reader is
	 tested on that exact output. Turn it on after one live restart plus a refresh cycle.
   - Tests: `TestPersistBlobPasswordHash*`, `TestDropRedundantPasswordHash`,
	 `TestValidateCloudKitRestoreStateUsesBlobHash`, Rust `persist_blob_matches_go_fixture`.
5. Encrypt credentials at rest with a key held in the macOS Keychain. Cover the `UserLoginMetadata` and
   session.json fields listed above. Reuse the AES-GCM code in `carddav_crypto.go`, but load the key from the
   Keychain instead of a key file. Check that a launchd agent can read the login keychain while the screen is
   locked (unverified). Effort: medium.
   Status: not implemented (as asked). Feasibility findings are below, under "Keychain key for credentials
   at rest".
6. Optional "no stored password" mode. Do not persist `hashed_password`, and accept a re-login when the PET
   expires (about 24 h; the exact lifetime is unverified). Touch `persist_data_from_spd`,
   `refresh_pet_token`, `login.go:796`. Effort: medium. Cost: CardDAV and CloudKit stop working until you log
   in again.
   Status: not implemented (not requested).
7. Add a config switch for live URL previews (`urlpreview.go`, `fetchURLPreview`) if leaking your IP to
   linked sites matters. Effort: small.
   Status: done, `1d0942c5` and `686ada1b`.
   - `network.url_previews`, default false. When off there is no `/preview_url` request and no page or
	 image GET. It gates five call sites: inbound live (`convertURLPreviewToBeeper`), the outbound
	 double-puppet edit (`addOutboundURLPreview`), outbound rich-link encoding
	 (`convertURLPreviewToIMessage`, which also scraped and was missing from the original list), CloudKit
	 backfill and chat.db backfill.
   - `url_previews_in_backfill` now applies only when `url_previews` is on
	 (`IMConfig.BackfillURLPreviews`).
   - Behaviour change: existing installs stop getting fetched previews after upgrading, until they set
	 `url_previews: true`.
   - Test: `TestURLPreviewsSwitchGatesLiveFetches` counts requests through a stub transport.
   - SSRF: `ogHTTPClient` had no address filter and followed up to 10 redirects (the `net/http` default).
	 It now dials through `rejectNonPublicDial`, which runs after DNS on every connection, so it covers
	 each redirect hop and DNS rebinding. It refuses loopback, RFC 1918/ULA, link-local (including
	 169.254.169.254), CGNAT 100.64/10 (the tailnet), unspecified, multicast, 0/8, NAT64 and
	 IPv4-mapped forms. `downloadURL` shares the client, so iCloud contact-avatar downloads are filtered
	 too. Tests: `TestIsPublicAddr`, `TestURLPreviewFetchRefusesLoopback`.
   - The homeserver's `/preview_url` is the homeserver's own SSRF problem; with `url_previews` off the
	 bridge no longer calls it.
   - Rich links without a fetch: yes. When the sender's device attached a rich link, rustpush delivers
	 it as `x-richlink/meta` (URL, canonical URL, title, description, image MIME) plus
	 `x-richlink/image` sideband attachments. `convertURLPreviewToBeeper` already turns those into a
	 `com.beeper.linkpreviews` entry with the uploaded image, and does no HTTP. That path is untouched
	 by the switch. Only messages without a rich link (SMS, old clients, a bare URL the sender's device
	 did not unfurl) lose their preview when `url_previews` is off. Outbound previews the Matrix client
	 supplies are still encoded into the iMessage without a fetch.

### Passwords-app client: removal cost (not removed)

- Verified statically, not at runtime: no Go caller. `rg PasswordsClient pkg/connector cmd pkg/cli` finds
  nothing outside the generated bindings. In `lib.rs`, `subsystem_state_path("passwords-state.plist")`
  appears only in `get_or_init_passwords_client`, which only the exported `get_passwords_client`
  reaches. So the plist is never written by this tree.
- Removal would need:
  - deleting `WrappedPasswordsClient` (about 105 lines, `lib.rs` 6316-6422), the `passwords_client`
	field and init, `get_or_init_passwords_client` and the export;
  - regenerating `rustpushgo.go` (89 references) and `rustpushgo.h` (31) with `uniffi-bindgen-go` plus
	`scripts/patch_bindings.py`. The generated Go checks each export's checksum at load time, so the
	bindings cannot simply be hand-edited. `uniffi-bindgen-go` is not installed here;
  - optionally dropping the two `passwords.rs` `rp_patch` lines in the Makefile and
	`tests/password_state_restore.rs`.
- Merge friction: high. The code came from upstream (`a7ce4572`, Cameron Aaron). `lib.rs` is already the
  worst conflict file, and the bindings get regenerated by every upstream sync. Each future merge would
  re-add or conflict on those hunks.
- Decision: left in place. It is dead code that writes nothing.

### Keychain key for credentials at rest (item 5, investigation only)

When it can be read:
- pitchfork `boot_start` and `bridge-all` run the bridge from a LaunchAgent. A LaunchAgent starts only
  when the user logs in, so the "after reboot, before login" case never runs the bridge at all. With
  FileVault on, auto-login is off, so a reboot means the bridge waits for a login whatever the key
  storage is.
- The login keychain is unlocked at login. Screen lock does not lock it unless "Lock when sleeping" or
  "Lock after N minutes" is set on that keychain; both are off by default. A LaunchAgent in the user's
  session can read it while the screen is locked. Unverified on this Mac: check with
  `security show-keychain-info ~/Library/Keychains/login.keychain-db`.
- `kSecAttrAccessibleAfterFirstUnlock` does not apply to the login (file-based) keychain. It is a
  data-protection-keychain attribute (`kSecUseDataProtectionKeychain`). On macOS that keychain needs a
  `keychain-access-groups` or application-identifier entitlement, which means a Developer ID or
  provisioning profile. The Makefile's ad-hoc `codesign --sign -` binary gets `errSecMissingEntitlement`
  (-34018). Without a paid signing identity, the login keychain is the only option.

Access control on the login keychain:
- Items carry a trusted-application list checked against the code signature's designated requirement.
  For an ad-hoc signature that requirement is the cdhash, which changes on every rebuild. Each rebuild
  would then raise a GUI "allow access" prompt that a headless agent cannot answer, and the read blocks
  or fails.
- Options:
  - (a) Sign with a stable self-signed code-signing certificate, so the requirement is
	"identifier com.lrhodin.corten-matrix and certificate leaf H". Rebuilds stay trusted, and other
	same-user processes are refused. This is the only option that also protects against other
	processes of the same user.
  - (b) Store the item with `security add-generic-password -A`. Any process of the user can read it
	without a prompt. That still protects backups, disk images and copied data dirs, but not malware
	running as the user.
  - (c) Read through `/usr/bin/security find-generic-password -w`. Then the trusted app is the
	Apple-signed `security` tool, which in practice equals (b), and the key passes through a pipe.

Which API:
- Go: `github.com/keybase/go-keychain` (cgo, Security.framework, `SecItemAdd`/`SecItemCopyMatching`).
  It fits, because the encrypting code (`carddav_crypto.go`, `identity_store.go`, `dbmeta.go`) is Go
  and the package already uses cgo. It is one new dependency.
- Rust: the `security-framework` crate (`passwords::get_generic_password`). It would add an FFI export
  and a bindings regeneration (see the Passwords section above), for no gain.
- `/usr/bin/security`: no dependency, but it has the ACL weakness in (c).

Where AES-GCM plugs in:
- `carddav_crypto.go` already has AES-256-GCM. Only `loadCardDAVKey`/`generateCardDAVKey` touch the
  0600 key file. Replace those two with a Keychain get-or-create (service `corten-matrix`, account
  `credentials-key`), keeping the file as a migration source.
- session.json: encrypt the marshalled JSON in `saveSessionStateLocked` before
  `writeSessionFileAtomically`. Decrypt in `loadSessionStateLocked`, plus `ListHandles` and
  `CheckSessionRestore`, which read the file directly. A magic prefix lets plaintext files load once and
  get re-encrypted on the next save.
- `UserLoginMetadata` (bridge DB): wrap the secret fields in a string type whose
  MarshalJSON/UnmarshalJSON encrypts and decrypts with an `enc:v1:` prefix. The fields are
  `AccountHashedPasswordHex`, `AccountPET`, `AccountPersistBlob`, `MmeDelegateJSON`, `IDSIdentity`,
  `APSState` and `IDSUsers`. The DB layer stays untouched, and unprefixed values read as legacy
  plaintext.
- Failure mode to design for: if the Keychain is unreadable (locked, prompt denied), the bridge must
  fail closed with a clear "keychain locked" bridge state, not treat the login as gone and wipe it.
- Each `bridge-all` account has its own data dir and needs its own Keychain account name, since both
  run as the same OS user.
- Effort stays medium: about 150 lines plus migration tests. The signing choice in (a) is the decision
  that decides how much it protects.

## Open decisions for the user

- Multi-user and per-user settings (a8713b33, a5f4a13f, c96e51a0) and the Docker flow: take them, or stay on
  lwittwer's one-process-per-account model?
- StatusKit topic style (5f85500c): ported with lwittwer's moon-in-title as the default. Nothing left to
  decide unless jason's "notice" style is wanted back.
- `url_previews` defaults to off, which changes behaviour for existing installs. Keep the default, or ship
  it on?
- `drop_redundant_password_hash`: turn it on after one live restart and a 12 h PET refresh succeed?
- Keychain encryption (item 5): it only helps against other same-user processes with a stable signing
  identity. Is a self-signed code-signing certificate acceptable?
- `keychain-join` (from 595ce7e9): port it? It asks for a device passcode in chat.
- Licence labelling of jason's hunks inside MPL-headed files.
- Credential hardening: items 2, 3, 4 (behind a flag) and 7 are done on `harden-credentials`; 5 and 6
  are open.
