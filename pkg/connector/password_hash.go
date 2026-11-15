// corten-matrix - A Matrix-iMessage puppeting bridge.
// Copyright (C) 2024 Ludvig Rhodin
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package connector

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
)

// accountPersistBlobVersion mirrors PERSIST_BLOB_VERSION in
// pkg/rustpushgo/src/lib.rs; restore_token_provider ignores any other version.
const accountPersistBlobVersion = 1

// effectivePasswordHashHex returns the hashed password restore and refresh
// will use: the legacy hex field when set, else the hash carried by the
// rustpush account blob, else "".
func effectivePasswordHashHex(hashHex, persistBlob, adsid string) string {
	if hashHex != "" {
		return hashHex
	}
	if hash := persistBlobPasswordHash(persistBlob, adsid); hash != nil {
		return hex.EncodeToString(hash)
	}
	return ""
}

// dropRedundantPasswordHash clears AccountHashedPasswordHex once the account
// blob is proven to carry the same hash in a form restore_token_provider will
// accept (same version, same ADSID). Returns true when the field was cleared.
func dropRedundantPasswordHash(meta *UserLoginMetadata) bool {
	if meta.AccountHashedPasswordHex == "" {
		return false
	}
	hash := persistBlobPasswordHash(meta.AccountPersistBlob, meta.AccountADSID)
	if hash == nil || hex.EncodeToString(hash) != strings.ToLower(meta.AccountHashedPasswordHex) {
		return false
	}
	meta.AccountHashedPasswordHex = ""
	return true
}

// persistBlobPasswordHash reads account.hashed_password from the base64
// binary plist written by encode_persist_blob. Returns nil unless the blob
// would be used by restore_token_provider for this account.
func persistBlobPasswordHash(persistBlob, adsid string) []byte {
	if persistBlob == "" || adsid == "" {
		return nil
	}
	raw, err := base64.StdEncoding.DecodeString(persistBlob)
	if err != nil {
		return nil
	}
	p, top, err := parseBplist(raw)
	if err != nil {
		return nil
	}
	if version, err := p.dictInt(top, "version"); err != nil || version != accountPersistBlobVersion {
		return nil
	}
	account, err := p.dictValue(top, "account")
	if err != nil {
		return nil
	}
	if blobADSID, err := p.dictBytes(account, "adsid", bplistASCII); err != nil || string(blobADSID) != adsid {
		return nil
	}
	hash, err := p.dictBytes(account, "hashed_password", bplistData)
	if err != nil || len(hash) == 0 {
		return nil
	}
	return hash
}

const (
	bplistInt   = 0x1
	bplistData  = 0x4
	bplistASCII = 0x5
	bplistDict  = 0xD
)

var errBplist = errors.New("malformed binary plist")

type bplist struct {
	data    []byte
	offsets []uint64
	refSize uint64
}

func parseBplist(data []byte) (*bplist, uint64, error) {
	if len(data) < 8+32 || string(data[:8]) != "bplist00" {
		return nil, 0, errBplist
	}
	trailer := data[len(data)-32:]
	offsetSize, refSize := uint64(trailer[6]), uint64(trailer[7])
	numObjects := binary.BigEndian.Uint64(trailer[8:16])
	top := binary.BigEndian.Uint64(trailer[16:24])
	tableStart := binary.BigEndian.Uint64(trailer[24:32])
	tableLimit := uint64(len(data) - 32)
	if offsetSize < 1 || offsetSize > 8 || refSize < 1 || refSize > 8 ||
		numObjects == 0 || numObjects > tableLimit || top >= numObjects ||
		tableStart > tableLimit || numObjects*offsetSize > tableLimit-tableStart {
		return nil, 0, errBplist
	}
	offsets := make([]uint64, numObjects)
	for i := range offsets {
		at := tableStart + uint64(i)*offsetSize
		offsets[i] = readBigEndian(data[at : at+offsetSize])
	}
	return &bplist{data: data, offsets: offsets, refSize: refSize}, top, nil
}

func readBigEndian(b []byte) uint64 {
	var v uint64
	for _, x := range b {
		v = v<<8 | uint64(x)
	}
	return v
}

// object returns an object's type nibble, its count (length, entry count or
// int byte-width exponent) and where its payload starts.
func (p *bplist) object(ref uint64) (kind byte, count, start uint64, err error) {
	if ref >= uint64(len(p.offsets)) || p.offsets[ref] >= uint64(len(p.data)) {
		return 0, 0, 0, errBplist
	}
	off := p.offsets[ref]
	marker := p.data[off]
	kind, count, start = marker>>4, uint64(marker&0x0F), off+1
	if kind == bplistInt || count != 0x0F {
		return kind, count, start, nil
	}
	if start >= uint64(len(p.data)) || p.data[start]>>4 != bplistInt {
		return 0, 0, 0, errBplist
	}
	width := uint64(1) << (p.data[start] & 0x0F)
	if width > 8 || start+1+width > uint64(len(p.data)) {
		return 0, 0, 0, errBplist
	}
	count = readBigEndian(p.data[start+1 : start+1+width])
	return kind, count, start + 1 + width, nil
}

func (p *bplist) bytesOf(ref uint64, want byte) ([]byte, error) {
	kind, length, start, err := p.object(ref)
	if err != nil || kind != want || length > uint64(len(p.data))-start {
		return nil, errBplist
	}
	return p.data[start : start+length], nil
}

func (p *bplist) intOf(ref uint64) (uint64, error) {
	kind, exp, start, err := p.object(ref)
	width := uint64(1) << exp
	if err != nil || kind != bplistInt || width > 8 || start+width > uint64(len(p.data)) {
		return 0, errBplist
	}
	return readBigEndian(p.data[start : start+width]), nil
}

func (p *bplist) dictValue(ref uint64, key string) (uint64, error) {
	kind, entries, start, err := p.object(ref)
	if err != nil || kind != bplistDict || entries > uint64(len(p.data)) ||
		2*entries*p.refSize > uint64(len(p.data))-start {
		return 0, errBplist
	}
	refAt := func(i uint64) uint64 {
		at := start + i*p.refSize
		return readBigEndian(p.data[at : at+p.refSize])
	}
	for i := uint64(0); i < entries; i++ {
		if k, err := p.bytesOf(refAt(i), bplistASCII); err == nil && string(k) == key {
			return refAt(entries + i), nil
		}
	}
	return 0, errBplist
}

func (p *bplist) dictBytes(ref uint64, key string, want byte) ([]byte, error) {
	value, err := p.dictValue(ref, key)
	if err != nil {
		return nil, err
	}
	return p.bytesOf(value, want)
}

func (p *bplist) dictInt(ref uint64, key string) (uint64, error) {
	value, err := p.dictValue(ref, key)
	if err != nil {
		return 0, err
	}
	return p.intOf(value)
}
