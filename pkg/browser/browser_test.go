package browser

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/khulnasoft-lab/go-goctl/v2/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GOCTL_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, "%v", os.Args[3:])
	os.Exit(0)
}

func TestBrowse(t *testing.T) {
	launcher := fmt.Sprintf("%q -test.run=TestHelperProcess -- chrome", os.Args[0])
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	b := Browser{launcher: launcher, stdout: stdout, stderr: stderr}
	err := b.browse("github.com", []string{"GOCTL_WANT_HELPER_PROCESS=1"})
	assert.NoError(t, err)
	assert.Equal(t, "[chrome github.com]", stdout.String())
	assert.Equal(t, "", stderr.String())
}

func TestResolveLauncher(t *testing.T) {
	tests := []struct {
		name         string
		env          map[string]string
		config       *config.Config
		wantLauncher string
	}{
		{
			name: "GOCTL_BROWSER set",
			env: map[string]string{
				"GOCTL_BROWSER": "GOCTL_BROWSER",
			},
			wantLauncher: "GOCTL_BROWSER",
		},
		{
			name:         "config browser set",
			config:       config.ReadFromString("browser: CONFIG_BROWSER"),
			wantLauncher: "CONFIG_BROWSER",
		},
		{
			name: "BROWSER set",
			env: map[string]string{
				"BROWSER": "BROWSER",
			},
			wantLauncher: "BROWSER",
		},
		{
			name: "GOCTL_BROWSER and config browser set",
			env: map[string]string{
				"GOCTL_BROWSER": "GOCTL_BROWSER",
			},
			config:       config.ReadFromString("browser: CONFIG_BROWSER"),
			wantLauncher: "GOCTL_BROWSER",
		},
		{
			name: "config browser and BROWSER set",
			env: map[string]string{
				"BROWSER": "BROWSER",
			},
			config:       config.ReadFromString("browser: CONFIG_BROWSER"),
			wantLauncher: "CONFIG_BROWSER",
		},
		{
			name: "GOCTL_BROWSER and BROWSER set",
			env: map[string]string{
				"BROWSER":    "BROWSER",
				"GOCTL_BROWSER": "GOCTL_BROWSER",
			},
			wantLauncher: "GOCTL_BROWSER",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != nil {
				for k, v := range tt.env {
					t.Setenv(k, v)
				}
			}
			if tt.config != nil {
				old := config.Read
				config.Read = func(_ *config.Config) (*config.Config, error) {
					return tt.config, nil
				}
				defer func() { config.Read = old }()
			}
			launcher := resolveLauncher()
			assert.Equal(t, tt.wantLauncher, launcher)
		})
	}
}
