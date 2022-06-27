package zipstd_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/batmac/zipstd"
)

func TestZipstd(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	testExamples(t)
}

func testExamples(t *testing.T) {
	t.Helper()
	n := 100
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		id := i
		go func() {
			t.Run(fmt.Sprintf("%v", id), testExample)
			wg.Done()
		}()
	}
	wg.Wait()
}

func testExample(t *testing.T) {
	t.Helper()
	// Try it out...
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	zipstd.Register(zw)
	// Create 1MB data
	tmp := make([]byte, 1<<20)
	/* 	for i := range tmp {
		tmp[i] = byte(i)
	} */
	nr, err := rand.Read(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if nr != len(tmp) {
		t.Fatal("tmp not fully written")
	}

	w, err := zw.CreateHeader(&zip.FileHeader{
		Name:   "file1.txt",
		Method: zipstd.ZipMethodWinZip,
	})
	if err != nil {
		t.Fatal(err)
	}
	n, err := w.Write(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatal("0 byte written")
	}
	// Another...
	w, err = zw.CreateHeader(&zip.FileHeader{
		Name:   "file2.txt",
		Method: zipstd.ZipMethodPKWare,
	})
	if err != nil {
		t.Fatal(err)
	}
	n, err = w.Write(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatal("0 byte written")
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	zipstd.Register(zr)
	for _, file := range zr.File {
		rc, err := file.Open()
		if err != nil {
			t.Fatal(err)
		}
		b, err := ioutil.ReadAll(rc)
		if err != nil {
			t.Fatal(err)
		}
		rc.Close()
		if bytes.Equal(b, tmp) {
			fmt.Println(file.Name, "ok", t.Name())
		} else {
			fmt.Println(file.Name, "mismatch", t.Name())
			t.Fail()
		}
	}
}
