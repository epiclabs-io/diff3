package diff3_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/epiclabs-io/diff3"
)

func compareReader(t *testing.T, a, b io.Reader) bool {
	abytes, err := io.ReadAll(a)
	if err != nil {
		t.Fatal(err)
	}
	bbytes, err := io.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.Equal(abytes, bbytes)
}

func TestDiff3(t *testing.T) {

	const testDir = "./testdata"
	const generate = false // set to true to genereate the expected files *-m.txt and *-error.txt
	files, err := os.ReadDir(testDir)
	if err != nil {
		log.Fatal(err)
	}

	tests := make(map[string]bool)

	for _, f := range files {
		if len(f.Name()) > 3 {
			tests[f.Name()[:3]] = true
		}
	}

	for k := range tests {
		func() {
			a, err := os.Open(fmt.Sprintf("%s/%s-a.txt", testDir, k))
			if err != nil {
				log.Fatal(err)
			}
			defer a.Close()
			b, err := os.Open(fmt.Sprintf("%s/%s-b.txt", testDir, k))
			if err != nil {
				log.Fatal(err)
			}
			defer b.Close()
			o, err := os.Open(fmt.Sprintf("%s/%s-o.txt", testDir, k))
			if err != nil {
				log.Fatal(err)
			}
			defer o.Close()
			var m io.ReadCloser
			if !generate {
				m, err = os.Open(fmt.Sprintf("%s/%s-m.txt", testDir, k))
				if err != nil {
					log.Fatal(err)
				}
				defer m.Close()
			}

			mr, mergeError := diff3.Merge(a, o, b, true, "A", "B")

			if generate {
				m, err := os.OpenFile(
					fmt.Sprintf("%s/%s-m.txt", testDir, k),
					os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
					0666,
				)
				if err != nil {
					t.Fatal(err)
				}
				defer m.Close()
				if mergeError == nil {
					_, err = io.Copy(m, mr.Result)
					if err != nil {
						t.Fatal(err)
					}
				}
				if mergeError != nil {
					err = os.WriteFile(fmt.Sprintf("%s/%s-error.txt", testDir, k), []byte(mergeError.Error()), 0666)
					if err != nil {
						t.Fatal(err)
					}
				}
			} else {
				var expectError = false
				if _, err := os.Stat(fmt.Sprintf("%s/%s-error.txt", testDir, k)); !os.IsNotExist(err) {
					expectError = true
				}
				if mergeError != nil && !expectError {
					t.Fatalf("Did not expect merge error: %s", err)
				}
				if !expectError && !compareReader(t, mr.Result, m) {
					t.Fatalf("Test #%s does not match expected result", k)
				}
			}
		}()
	}
}

func TestMergeWithOptionsDefaultsToLCS(t *testing.T) {
	a := "one\ntwo modified\nthree"
	o := "one\ntwo\nthree"
	b := "one\ntwo\nthree modified"

	defaultResult, err := diff3.Merge(strings.NewReader(a), strings.NewReader(o), strings.NewReader(b), true, "A", "B")
	if err != nil {
		t.Fatal(err)
	}
	defaultBytes, err := io.ReadAll(defaultResult.Result)
	if err != nil {
		t.Fatal(err)
	}

	optionsResult, err := diff3.MergeWithOptions(strings.NewReader(a), strings.NewReader(o), strings.NewReader(b), diff3.MergeOptions{
		Algorithm:             diff3.DiffAlgorithmLCS,
		ExcludeFalseConflicts: true,
		Detailed:              true,
		LabelA:                "A",
		LabelB:                "B",
	})
	if err != nil {
		t.Fatal(err)
	}
	optionsBytes, err := io.ReadAll(optionsResult.Result)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(defaultBytes, optionsBytes) {
		t.Fatalf("default merge and LCS options merge differ")
	}
}

func TestMergeWithOptionsMyers(t *testing.T) {
	a := "one\ntwo modified\nthree"
	o := "one\ntwo\nthree"
	b := "one\ntwo\nthree"

	result, err := diff3.MergeWithOptions(strings.NewReader(a), strings.NewReader(o), strings.NewReader(b), diff3.MergeOptions{
		Algorithm:             diff3.DiffAlgorithmMyers,
		ExcludeFalseConflicts: true,
		Detailed:              true,
		LabelA:                "A",
		LabelB:                "B",
	})
	if err != nil {
		t.Fatal(err)
	}
	merged, err := io.ReadAll(result.Result)
	if err != nil {
		t.Fatal(err)
	}

	if result.Conflicts {
		t.Fatalf("expected clean merge, got conflict in %q", string(merged))
	}
	if got, want := string(merged), "one\ntwo modified\nthree"; got != want {
		t.Fatalf("unexpected Myers merge output: got %q, want %q", got, want)
	}
}
