package helmexec

import "go.uber.org/zap"

type nativeExecer struct {
	execer
}

func (self *nativeExecer) IsHelm3() bool {
	return self.version.Major() >= 3
}

func NewNativeExec(helmBinary string, options HelmExecOptions, logger *zap.SugaredLogger, kubeconfig string, kubeContext string, runner Runner) *nativeExecer {
	// TODO: proper error handling
	version, err := GetHelmVersion(helmBinary, runner)
	if err != nil {
		panic(err)
	}

	if version.Prerelease() != "" {
		logger.Warnf("Helm version %s is a pre-release version. This may cause problems when deploying Helm charts.\n", version)
		*version, _ = version.SetPrerelease("")
	}

	return &nativeExecer{
		execer: execer{
			helmBinary:       helmBinary,
			options:          options,
			version:          version,
			logger:           logger,
			kubeconfig:       kubeconfig,
			kubeContext:      kubeContext,
			runner:           runner,
			decryptedSecrets: make(map[string]*decryptedSecret),
		},
	}
}
