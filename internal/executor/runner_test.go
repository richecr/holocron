package executor

import (
	"runtime"
	"testing"

	"github.com/richecr/holocron/internal/config"
)

func TestResolveCommandPriority(t *testing.T) {
	goos := runtime.GOOS

	tests := []struct {
		name    string
		script  config.ScriptConfig
		want    string
		wantErr bool
	}{
		{
			name: "platform specific wins over generic",
			script: config.ScriptConfig{
				Command: "generic-cmd",
				Platforms: map[string]string{
					goos: "platform-cmd",
				},
			},
			want: "platform-cmd",
		},
		{
			name: "generic fallback when no platform entry",
			script: config.ScriptConfig{
				Command: "generic-cmd",
				Platforms: map[string]string{
					"some-other-os": "other-cmd",
				},
			},
			want: "generic-cmd",
		},
		{
			name: "empty platform entry falls back to generic",
			script: config.ScriptConfig{
				Command: "generic-cmd",
				Platforms: map[string]string{
					goos: "",
				},
			},
			want: "generic-cmd",
		},
		{
			name: "error when nothing matches",
			script: config.ScriptConfig{
				Platforms: map[string]string{
					"some-other-os": "other-cmd",
				},
			},
			wantErr: true,
		},
		{
			name:    "error when completely empty",
			script:  config.ScriptConfig{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveCommand(tt.script)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ResolveCommand() = %q, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ResolveCommand() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("ResolveCommand() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSelectShell(t *testing.T) {
	tests := []struct {
		name     string
		settings config.Settings
		goos     string
		want     string
		wantFlag string
	}{
		{
			name:     "linux default is /bin/sh -c",
			settings: config.Settings{DefaultShell: map[string]string{}},
			goos:     "linux",
			want:     "/bin/sh",
			wantFlag: "-c",
		},
		{
			name:     "darwin default is /bin/sh -c",
			settings: config.Settings{DefaultShell: map[string]string{}},
			goos:     "darwin",
			want:     "/bin/sh",
			wantFlag: "-c",
		},
		{
			name:     "windows default is powershell.exe",
			settings: config.Settings{DefaultShell: map[string]string{}},
			goos:     "windows",
			want:     "powershell.exe",
			wantFlag: "-NoProfile",
		},
		{
			name: "config overrides linux default with bash",
			settings: config.Settings{
				DefaultShell: map[string]string{"linux": "bash"},
			},
			goos:     "linux",
			want:     "bash",
			wantFlag: "-c",
		},
		{
			name: "config overrides windows default with cmd",
			settings: config.Settings{
				DefaultShell: map[string]string{"windows": "cmd"},
			},
			goos:     "windows",
			want:     "cmd.exe",
			wantFlag: "/c",
		},
		{
			name: "config powershell alias keeps -NoProfile",
			settings: config.Settings{
				DefaultShell: map[string]string{"windows": "powershell"},
			},
			goos:     "windows",
			want:     "powershell.exe",
			wantFlag: "-NoProfile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotFlag := selectShell(tt.settings, tt.goos)
			if got != tt.want {
				t.Errorf("selectShell() shell = %q, want %q", got, tt.want)
			}
			if gotFlag != tt.wantFlag {
				t.Errorf("selectShell() flag = %q, want %q", gotFlag, tt.wantFlag)
			}
		})
	}
}

func TestBuildInvocation(t *testing.T) {
	tests := []struct {
		shell   string
		flag    string
		command string
		want    []string
	}{
		{"/bin/sh", "-c", "echo hi", []string{"-c", "echo hi"}},
		{"bash", "-c", "echo hi", []string{"-c", "echo hi"}},
		{"powershell.exe", "-NoProfile", "echo hi", []string{"-NoProfile", "-Command", "echo hi"}},
		{"cmd.exe", "/c", "echo hi", []string{"/c", "echo hi"}},
	}

	for _, tt := range tests {
		got := buildInvocation(tt.shell, tt.flag, tt.command)
		if len(got) != len(tt.want) {
			t.Fatalf("buildInvocation(%q) = %v, want %v", tt.shell, got, tt.want)
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("buildInvocation(%q)[%d] = %q, want %q", tt.shell, i, got[i], tt.want[i])
			}
		}
	}
}
