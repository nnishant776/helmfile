package helmexec

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/helmfile/helmfile/pkg/envvar"
	"go.uber.org/zap"
	helmcmd "helm.sh/helm/v4/pkg/cmd"
)

// NativeRunner implemention for shell commands
type NativeRunner struct {
	Dir string

	StripArgsValuesOnExitError bool

	Logger *zap.SugaredLogger
	Ctx    context.Context
}

// Execute helm as a native command
func (self *NativeRunner) Execute(cmd string, args []string, env map[string]string, enableLiveOutput bool) ([]byte, error) {
	logWriterGenerators := []*logWriterGenerator{
		&logWriterGenerator{
			log: self.Logger,
		},
	}

	id := ""
	if os.Getenv(envvar.DisableRunnerUniqueID) == "" {
		id = newExecutionID()
	}

	var logWriters []io.Writer
	for _, g := range logWriterGenerators {
		var logPrefix string
		if id == "" {
			logPrefix = fmt.Sprintf("%s> ", "helm")
		} else {
			logPrefix = fmt.Sprintf("%s:%s> ", "helm", id)
		}
		logWriters = append(logWriters, g.Writer(logPrefix))
	}

	stdout := bytes.Buffer{}
	stdoutWriters := io.MultiWriter(append(logWriters, &stdout)...)
	c, err := helmcmd.NewRootCmd(stdoutWriters, args, func(b bool) {})
	if err != nil {
		return nil, err
	}

	c.SetArgs(args)
	err = c.ExecuteContext(self.Ctx)
	if err != nil {
		return nil, err
	}

	return stdout.Bytes(), nil
}

// Execute a shell command
func (self *NativeRunner) ExecuteStdIn(cmd string, args []string, env map[string]string, stdin io.Reader) ([]byte, error) {
	_, err := io.Copy(os.Stdin, stdin)
	if err != nil {
		return nil, err
	}

	return self.Execute(cmd, args, env, true)
}
