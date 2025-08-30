package runtime

import (
	"fmt"
	"os"
	"strconv"

	"github.com/helmfile/helmfile/pkg/envvar"
)

var (
	// GoYamlV3 is set to true in order to let Helmfile use
	// go.yaml.in/yaml/v3 instead of go.yaml.in/yaml/v2.
	// It's false by default in Helmfile v0.x and true in Helmfile v1.x.
	GoYamlV3 bool

	// NativeHelm indicates whether helmfile should explicitly look
	// for a helm binary and run it, or it should directly use the
	// currently experimental helm/v4 API to directly pass the command
	// arguments to the cobra command.
	// This will be set to false by default and can be controlled by the
	// environment variable NATIVE_HELM
	NativeHelm bool
)

func Info() string {
	yamlLib := "go.yaml.in/yaml/v2"
	if GoYamlV3 {
		yamlLib = "go.yaml.in/yaml/v3"
	}

	return fmt.Sprintf("YAML library = %v", yamlLib)
}

func init() {
	// You can switch the YAML library at runtime via an envvar:
	switch os.Getenv(envvar.GoYamlV3) {
	case "true":
		GoYamlV3 = true
	case "false":
		GoYamlV3 = false
	default:
		GoYamlV3 = true
	}

	if v, err := strconv.ParseBool(os.Getenv("NATIVE_HELM")); err == nil && v {
		NativeHelm = true
	}
}
