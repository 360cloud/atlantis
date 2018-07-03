package runtime

import (
	"github.com/hashicorp/go-version"
	"github.com/runatlantis/atlantis/server/events/models"
)

// InitStep runs `terraform init`.
type InitStepRunner struct {
	TerraformExecutor TerraformExec
	DefaultTFVersion  *version.Version
}

// nolint: unparam
func (i *InitStepRunner) Run(ctx models.ProjectCommandContext, extraArgs []string, path string) (string, error) {
	tfVersion := i.DefaultTFVersion
	if ctx.ProjectConfig != nil && ctx.ProjectConfig.TerraformVersion != nil {
		tfVersion = ctx.ProjectConfig.TerraformVersion
	}
	// If we're running < 0.9 we have to use `terraform get` instead of `init`.
	if MustConstraint("< 0.9.0").Check(tfVersion) {
		ctx.Log.Info("running terraform version < 0.9.0 so will use `get` instead of `init`", "version", tfVersion)
		terraformGetCmd := append([]string{"get", "-no-color"}, extraArgs...)
		_, err := i.TerraformExecutor.RunCommandWithVersion(ctx.Log, path, terraformGetCmd, tfVersion, ctx.Workspace)
		return "", err
	} else {
		_, err := i.TerraformExecutor.RunCommandWithVersion(ctx.Log, path, append([]string{"init", "-no-color"}, extraArgs...), tfVersion, ctx.Workspace)
		return "", err
	}
}
