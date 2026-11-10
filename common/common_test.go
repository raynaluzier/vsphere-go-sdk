package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInsecureSkipTLSVerificationRequiresExplicitTrue(t *testing.T) {
	t.Setenv("VSPHERE_INSECURE_SKIP_TLS_VERIFY", "")
	if InsecureSkipTLSVerification() {
		t.Fatal("TLS verification should remain enabled when the compatibility switch is unset")
	}

	t.Setenv("VSPHERE_INSECURE_SKIP_TLS_VERIFY", "true")
	if !InsecureSkipTLSVerification() {
		t.Fatal("expected explicit true to enable the compatibility switch")
	}
}

func TestVcenterAuthReturnsEmptyTokenForInvalidServerURL(t *testing.T) {
	if got := VcenterAuth("user", "password", ":"); got != "" {
		t.Fatalf("token = %q, want empty token for invalid server URL", got)
	}
}

func TestCheckDirAndCreateReturnsFailureForFilePath(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(filePath, []byte("file"), 0600); err != nil {
		t.Fatal(err)
	}

	if got := CheckDirAndCreate(filePath); got != "Failed directory check" {
		t.Fatalf("result = %q, want failed directory check", got)
	}
}
