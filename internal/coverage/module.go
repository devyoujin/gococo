package coverage

import (
	"fmt"
	"os"
	"os/exec"
)

type module struct {
	path         string
	name         string
}

func (module *module) generateCoverage(dir string) error {
	cmd := exec.Command("go", "test", "-cover", "./...", "-test.gocoverdir="+dir)
	cmd.Dir = module.path
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run tests in %s: %w", module.path, err)
	}
	return nil
}