package job

import (
	"github.com/beoboo/job-scheduler/library/helpers"
	"github.com/beoboo/job-scheduler/library/log"
	"os"
	"os/exec"
)

type Child struct {
}

func (c *Child) Run(jobId, executable string, args ...string) {
	log.Debugf("[%s] Executing \"%s\"\n", jobId, helpers.FormatCmdLine(executable, args...))
	cmd := exec.Command(executable, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Unexpected: %s\n", err)
	}
}
