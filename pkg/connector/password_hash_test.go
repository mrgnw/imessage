package connector

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

// Written by rustpush's encode_persist_blob; pinned by
// persist_blob_matches_go_fixture in pkg/rustpushgo/src/lib.rs.
const rustPersistBlobFixture = "YnBsaXN0MDDSAQIDBFd2ZXJzaW9uV2FjY291bnQQAdYFBgcICQoLDA0ODwtYdXNlcm5hbWVfEA9oYXNoZWRfcGFzc3dvcmRWdG9rZW5zVWFkc2lkVGRzaWRWYWNuYW1lXxAUdXNlckBleGFtcGxlLmludmFsaWRPECAAAQIDBAUGBwgJCgsMDQ4PEBESExQVFhcYGRobHB0eH9EQEV8QFWNvbS5hcHBsZS5ncy5pZG1zLnBldNISExQVVXRva2VuWmV4cGlyYXRpb25TcGV00hYXGBhfEBBzZWNzX3NpbmNlX2Vwb2NoXxARbmFub3Nfc2luY2VfZXBvY2gQAF4wMDAxMTEtMjItYWFhYRJJlgLSAAgADQAVAB0AHwAsADUARwBOAFQAWQBgAHcAmgD9AQwAnQC1ALoAwADLAM8A1ADnAPsAAAAAAAACAQAAAAAAAAAZAAAAAAAAAAAAAAAAAAABEQ=="

// Same shape written by Python's plistlib (1-byte offsets, extra provider key).
const pythonPersistBlobFixture = "YnBsaXN0MDDTAQIDBAUXV3ZlcnNpb25XYWNjb3VudFhwcm92aWRlchAB1gYHCAkKCwwNDhUWDFh1c2VybmFtZV8QD2hhc2hlZF9wYXNzd29yZFZ0b2tlbnNVYWRzaWRUZHNpZFZhY25hbWVfEBR1c2VyQGV4YW1wbGUuaW52YWxpZE8QIAABAgMEBQYHCAkKCwwNDg8QERITFBUWFxgZGhscHR4f0Q8QXxAVY29tLmFwcGxlLmdzLmlkbXMucGV00hESExRVdG9rZW5aZXhwaXJhdGlvblNwZXQQAF4wMDAxMTEtMjItYWFhYRJJlgLS0RgZVXN0YXRlVm9wYXF1ZQAIAA8AFwAfACgAKgA3AEAAUgBZAF8AZABrAIIApQCoAMAAxQDLANYA2gDcAOsA8ADzAPkAAAAAAAACAQAAAAAAAAAaAAAAAAAAAAAAAAAAAAABAA=="

const fixtureADSID = "000111-22-aaaa"

var fixtureHashHex = func() string {
	hash := make([]byte, 32)
	for i := range hash {
		hash[i] = byte(i)
	}
	return hex.EncodeToString(hash)
}()

func TestPersistBlobPasswordHashReadsRustBlob(t *testing.T) {
	for name, blob := range map[string]string{"rust": rustPersistBlobFixture, "python": pythonPersistBlobFixture} {
		if got := hex.EncodeToString(persistBlobPasswordHash(blob, fixtureADSID)); got != fixtureHashHex {
			t.Fatalf("%s fixture: hash = %q, want %q", name, got, fixtureHashHex)
		}
	}
}

func TestPersistBlobPasswordHashRejectsUnusableBlobs(t *testing.T) {
	raw, err := base64.StdEncoding.DecodeString(rustPersistBlobFixture)
	if err != nil {
		t.Fatal(err)
	}
	truncated := base64.StdEncoding.EncodeToString(raw[:len(raw)-10])
	wrongVersion := append([]byte(nil), raw...)
	wrongVersion[strings.Index(string(raw), "account")+len("account")+1] = 2
	cases := map[string]struct{ blob, adsid string }{
		"empty blob":      {"", fixtureADSID},
		"empty adsid":     {rustPersistBlobFixture, ""},
		"other account":   {rustPersistBlobFixture, "other-adsid"},
		"not base64":      {"%%%", fixtureADSID},
		"not a plist":     {base64.StdEncoding.EncodeToString([]byte("bplist00 but not really, padded to forty bytes")), fixtureADSID},
		"truncated":       {truncated, fixtureADSID},
		"unknown version": {base64.StdEncoding.EncodeToString(wrongVersion), fixtureADSID},
	}
	for name, tc := range cases {
		if got := persistBlobPasswordHash(tc.blob, tc.adsid); got != nil {
			t.Fatalf("%s: got hash %x, want nil", name, got)
		}
	}
}

func TestDropRedundantPasswordHash(t *testing.T) {
	t.Run("blob carries the same hash", func(t *testing.T) {
		meta := &UserLoginMetadata{
			AccountHashedPasswordHex: strings.ToUpper(fixtureHashHex),
			AccountPersistBlob:       rustPersistBlobFixture,
			AccountADSID:             fixtureADSID,
		}
		if !dropRedundantPasswordHash(meta) || meta.AccountHashedPasswordHex != "" {
			t.Fatalf("hex copy was not dropped: %q", meta.AccountHashedPasswordHex)
		}
		if got := effectivePasswordHashHex(meta.AccountHashedPasswordHex, meta.AccountPersistBlob, meta.AccountADSID); got != fixtureHashHex {
			t.Fatalf("restore would use %q, want the blob hash", got)
		}
	})
	keptCases := map[string]*UserLoginMetadata{
		"no blob":         {AccountHashedPasswordHex: fixtureHashHex, AccountADSID: fixtureADSID},
		"different hash":  {AccountHashedPasswordHex: "aabb", AccountPersistBlob: rustPersistBlobFixture, AccountADSID: fixtureADSID},
		"blob of another": {AccountHashedPasswordHex: fixtureHashHex, AccountPersistBlob: rustPersistBlobFixture, AccountADSID: "other"},
	}
	for name, meta := range keptCases {
		t.Run(name, func(t *testing.T) {
			before := meta.AccountHashedPasswordHex
			if dropRedundantPasswordHash(meta) || meta.AccountHashedPasswordHex != before {
				t.Fatalf("hex copy must be kept, got %q", meta.AccountHashedPasswordHex)
			}
		})
	}
}

func TestValidateCloudKitRestoreStateUsesBlobHash(t *testing.T) {
	spd := `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0"><dict>
<key>DsPrsId</key><integer>1234567890</integer>
<key>adsid</key><string>000111-22-aaaa</string>
</dict></plist>`
	delegate := `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0"><dict>
<key>tokens</key><dict><key>com.apple.mobileme</key><string>synthetic-token</string></dict>
<key>config</key><dict>
<key>com.apple.Dataclass.KeychainSync</key><dict>
<key>escrowProxyUrl</key><string>https://escrow.example.invalid</string>
</dict></dict>
</dict></plist>`
	state := PersistedSessionState{
		AccountUsername:    "user@example.invalid",
		AccountPET:         "pet",
		AccountADSID:       fixtureADSID,
		AccountSPDBase64:   base64.StdEncoding.EncodeToString([]byte(spd)),
		MmeDelegateJSON:    delegate,
		AccountPersistBlob: rustPersistBlobFixture,
	}
	if err := validateCloudKitRestoreState(state); err != nil {
		t.Fatalf("hash held only in the blob was rejected: %v", err)
	}
	state.AccountPersistBlob = ""
	if err := validateCloudKitRestoreState(state); err == nil {
		t.Fatal("state with no hash anywhere was accepted")
	}
}
