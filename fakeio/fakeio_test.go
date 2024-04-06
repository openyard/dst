package fakeio_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/openyard/dst/fakeio"
	"github.com/stretchr/testify/assert"
)

const (
	Byte = 1
	KB   = 1024 * Byte
	MB   = 1024 * KB
	GB   = 1024 * MB
)

func TestFakeReadOnlyFile_HashAndSend_Small(t *testing.T) {
	start := time.Now()
	f := fakeio.NewFakeReadOnlyFile(KB)
	b, sha := hashAndSend(t, &f, KB)
	t.Logf("took %v to hash!", time.Since(start))

	assert.Equal(t, `abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghij`, string(b))
	assert.Equal(t, "dba4a6315b76548b7a4dd079ef6aa29a7b34fa8b92c11668473441715c5f0af5", sha)
}

func TestFakeReadOnlyFile_HashAndSend_Big(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping [TestFakeReadOnlyFile_HashAndSend_Big] with testing.Short()")
	}
	start := time.Now()
	f := fakeio.NewFakeReadOnlyFile(GB * 100) // needs around 10min
	_, sha := hashAndSend(t, &f, GB)
	t.Logf("took %v to hash!", time.Since(start))

	assert.Equal(t, "162e7af51048e81bec4a0b20851f181ff83b999f718baa2d839aa2412ab34406", sha)
}

func hashAndSend(t *testing.T, f *fakeio.FakeReadOnlyFile, size int) ([]byte, string) {
	t.Helper()
	w := sha256.New()

	//any reads from tee will read from r and write to w
	tee := io.TeeReader(f, w)

	b := sendReader(t, tee, size)
	sha := hex.EncodeToString(w.Sum(nil))
	return b, sha
}

func sendReader(t *testing.T, data io.Reader, size int) []byte {
	t.Helper()
	buff := make([]byte, size)
	for {
		_, err := data.Read(buff)
		if err == io.EOF {
			break
		}
		fmt.Sprintln(string(buff))
	}
	return buff
}
