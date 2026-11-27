package service

import (
	"reflect"
	"testing"
)

func TestAppVersionHelpers(t *testing.T) {
	if !isNewerAppVersion("demo", "1.2.0", "1.1.9") {
		t.Fatal("expected a newer numeric version")
	}
	if isNewerAppVersion("demo", "latest", "1.1.9") {
		t.Fatal("non-numeric version should not be ordered")
	}
	if !isCrossAppVersion("demo", "1.9.0", "2.0.0") {
		t.Fatal("expected a cross-major upgrade")
	}
	if isCrossAppVersion("demo", "2.0.0", "1.9.0") {
		t.Fatal("a downgrade should not be a cross-version upgrade")
	}
	got := sortAppVersions("demo", []string{"1.1", "1.3", "1.2"})
	want := []string{"1.3", "1.2", "1.1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sortAppVersions() = %v, want %v", got, want)
	}
}
