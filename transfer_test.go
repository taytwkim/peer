package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
)

func TestTransferResponseHeaderDoesNotLeakDelimiterIntoBody(t *testing.T) {
	var stream bytes.Buffer
	body := []byte("{manifest-json-starts-like-this}")

	if err := writeTransferResponseHeader(&stream, TransferResponse{
		Kind:     string(ObjectManifest),
		Filesize: int64(len(body)),
		Filename: "manifest.json",
	}); err != nil {
		t.Fatal(err)
	}
	stream.Write(body)

	var resp TransferResponse
	bodyReader, err := readTransferResponseHeader(&stream, &resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Kind != string(ObjectManifest) {
		t.Fatalf("kind = %q, want %q", resp.Kind, ObjectManifest)
	}

	got, err := io.ReadAll(io.LimitReader(bodyReader, resp.Filesize))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, body) {
		t.Fatalf("body = %q, want %q", got, body)
	}
}

func TestFilterSelfProviderCandidates(t *testing.T) {
	self := peer.ID("self")
	other := peer.ID("other")
	providers := []peer.AddrInfo{
		{ID: self},
		{ID: other},
		{ID: self},
	}

	got := filterSelfProviderCandidates(providers, self)
	if len(got) != 1 {
		t.Fatalf("filtered providers = %v, want only other provider", got)
	}
	if got[0].ID != other {
		t.Fatalf("filtered provider = %s, want %s", got[0].ID, other)
	}
}
