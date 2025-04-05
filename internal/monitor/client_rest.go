package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/chlp/ui/internal/model"
	"net/http"
)

func getRestInfo(address string) (*model.DeviceInfo, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/info", address))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}
	var info model.DeviceInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}
