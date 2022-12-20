package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
)

// RunCommandProvisioner returns a provisioner for the `docker run` command, which runs the specified provisioner
// and then forwards its provision output to the run command:
// 1. Environment variable get forwarded by setting the environment variable names as `-e` flags
// 2. Files get forwarded by mounting the temp dir as a volume and replacing all occurences of the temp dir path
// with the mounted dir path (i.e. `/.op/secrets`). Mounting fixed paths is not supported.
// 3. Args get forwarded by appending them to the run command
func RunCommandProvisioner(forProvisioner sdk.Provisioner) sdk.Provisioner {
	return runCommandProvisioner{forProvisioner: forProvisioner}
}

type runCommandProvisioner struct {
	sdk.Provisioner

	forProvisioner sdk.Provisioner
}

const SecretsVolumeMountDestination = "/.op/secrets"

func (p runCommandProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	p.forProvisioner.Provision(ctx, in, out)

	mountSecretsDir := false
	for path := range out.Files {
		if strings.HasPrefix(path, in.TempDir) {
			mountSecretsDir = true
			break
		}
	}

	if mountSecretsDir {
		for envVarName, value := range out.Environment {
			if strings.HasPrefix(value, in.TempDir) {
				valueWithVolumePath := strings.Replace(value, in.TempDir, SecretsVolumeMountDestination, 1)
				out.AddEnvVar(envVarName, valueWithVolumePath)
			}
		}

		for i, arg := range out.CommandLine {
			if strings.Contains(arg, in.TempDir) {
				arg = strings.ReplaceAll(arg, in.TempDir, SecretsVolumeMountDestination)
			}
			out.CommandLine[i] = arg
		}

		out.AddArgs("-v", fmt.Sprintf("%s:%s", in.TempDir, SecretsVolumeMountDestination))
	}

	// Forward all environment variable down to `docker run`, using the `-e` flag with the env var name
	for envVarName := range out.Environment {
		out.AddArgs("-e", envVarName)
	}
}

func (p runCommandProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	p.forProvisioner.Deprovision(ctx, in, out)
}

func (p runCommandProvisioner) Description() string {
	return p.forProvisioner.Description()
}
