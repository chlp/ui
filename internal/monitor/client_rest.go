package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/chlp/ui/internal/device"
	"net/http"
	"time"
)

const restClientTimeout = 2 * time.Second

func getRestInfo(address string) (*device.Info, error) {
	client := &http.Client{
		Timeout: restClientTimeout,
	}
	resp, err := client.Get(fmt.Sprintf("http://%s/v1/info", address))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}
	var info device.Info
	if err = json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}
