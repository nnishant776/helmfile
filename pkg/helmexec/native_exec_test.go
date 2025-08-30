package helmexec

import (
	"os"
	"testing"
)

func Test_IsHelm3Native(t *testing.T) {
	helm2Runner := mockRunner{output: []byte("Client: v2.16.0+ge13bc94\n")}
	helm := NewNativeExec("helm", HelmExecOptions{}, NewLogger(os.Stdout, "info"), "", "dev", &helm2Runner)
	if helm.IsHelm3() {
		t.Error("helmexec.IsHelm3() - Detected Helm 3 with Helm 2 version")
	}

	helm3Runner := mockRunner{output: []byte("v3.0.0+ge29ce2a\n")}
	helm = NewNativeExec("helm", HelmExecOptions{}, NewLogger(os.Stdout, "info"), "", "dev", &helm3Runner)
	if !helm.IsHelm3() {
		t.Error("helmexec.IsHelm3() - Failed to detect Helm 3")
	}

	helm4Runner := mockRunner{output: []byte("v4.0.0+ge29ce2a\n")}
	helm = NewNativeExec("helm", HelmExecOptions{}, NewLogger(os.Stdout, "info"), "", "dev", &helm4Runner)
	if !helm.IsHelm3() {
		t.Error("helmexec.IsHelm3() - Failed to detect Helm 3")
	}
}
