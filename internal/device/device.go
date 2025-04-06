package device

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/chlp/ui/pkg/exec"
)

func GetLocalDevice(info *Info, checksumCmd string, checksumEmulate bool) *Info {
	if info == nil {
		return nil
	}

	info.Checksum = exec.GetStringFromCmd(checksumCmd)
	if info.Checksum == "" && checksumEmulate {
		hash := sha256.Sum256([]byte(info.ID + info.HardwareVersion + info.SoftwareVersion + info.FirmwareVersion))
		info.Checksum = hex.EncodeToString(hash[:])
	}

	return info
}
