package exec

import (
	"github.com/chlp/ui/pkg/logger"
	"os/exec"
	"strings"
)

func GetStringFromCmd(cmd string) string {
	parts := strings.Fields(strings.TrimSpace(cmd))
	if len(parts) == 0 {
		return ""
	}

	output, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		logger.Printf("ExecGetStringFromCmd: err %v", err)
		return ""
	}

	return strings.TrimSpace(string(output))
}
