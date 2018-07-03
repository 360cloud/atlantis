package runtime_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/runatlantis/atlantis/server/events/models"
	"github.com/runatlantis/atlantis/server/events/runtime"
	"github.com/runatlantis/atlantis/server/events/yaml/valid"
	. "github.com/runatlantis/atlantis/testing"
	log "gopkg.in/inconshreveable/log15.v2"
)

func TestRunStepRunner_Run(t *testing.T) {
	cases := []struct {
		Command string
		ExpOut  string
		ExpErr  string
	}{
		{
			Command: "",
			ExpErr:  "no commands for run step",
		},
		{
			Command: "echo hi",
			ExpOut:  "hi\n",
		},
		{
			Command: "echo hi >> file && cat file",
			ExpOut:  "hi\n",
		},
		{
			Command: "lkjlkj",
			ExpErr:  "exit status 127: running \"lkjlkj\" in",
		},
		{
			Command: "echo workspace=$WORKSPACE version=$ATLANTIS_TERRAFORM_VERSION dir=$DIR",
			ExpOut:  "workspace=myworkspace version=0.11.0 dir=$DIR\n",
		},
	}

	projVersion, err := version.NewVersion("v0.11.0")
	Ok(t, err)
	defaultVersion, _ := version.NewVersion("0.8")
	r := runtime.RunStepRunner{
		DefaultTFVersion: defaultVersion,
	}
	ctx := models.ProjectCommandContext{
		Log:        log.New(),
		Workspace:  "myworkspace",
		RepoRelDir: "mydir",
		ProjectConfig: &valid.Project{
			TerraformVersion: projVersion,
			Workspace:        "myworkspace",
			Dir:              "mydir",
		},
	}
	for _, c := range cases {
		t.Run(c.Command, func(t *testing.T) {
			tmpDir, cleanup := TempDir(t)
			defer cleanup()
			var split []string
			if c.Command != "" {
				split = strings.Split(c.Command, " ")
			}
			out, err := r.Run(ctx, split, tmpDir)
			if c.ExpErr != "" {
				ErrContains(t, c.ExpErr, err)
				return
			}
			Ok(t, err)
			expOut := strings.Replace(c.ExpOut, "dir=$DIR", "dir="+tmpDir, -1)
			Equals(t, expOut, out)
		})
	}
}
