
package web

import (
	"fmt"
	"testing"
)

func assertEqual[K comparable](t *testing.T, want K, got K) {
	t.Helper()

	if want != got {
		t.Fatal(fmt.Printf("expected: '%v', got: '%v'", want, got))
	}
}

func TestEncodeTitle(t *testing.T) {
	want := "I+want+pizzqa+!"

	got := encodeTitle([]byte("  I want pizzqa !   "))

	assertEqual(t, want, string(got))
}

func TestDecodeTitle(t *testing.T) {
	want := "I want pizzqa !"

	got := decodeTitle([]byte("I+want+pizzqa+!"))

	assertEqual(t, want, string(got))
}

func TestAddTitleReferences(t *testing.T) {
	arg := []byte("article: [[hi me]]")
	want := `article: <a href="/view/hi+me" class="reference">hi me</a>`

	got := addTitleReferences(arg)

	assertEqual(t, want, string(got))
}
