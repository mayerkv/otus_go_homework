package main

import "testing"

func TestRunCmd(t *testing.T) {
	type args struct {
		cmd []string
		env Environment
	}
	tests := []struct {
		name           string
		args           args
		wantReturnCode int
	}{
		{
			name: "with returns 0",
			args: args{
				cmd: []string{"ls", "-la"},
				env: nil,
			},
			wantReturnCode: 0,
		},
		{
			name: "with returns 1",
			args: args{
				cmd: []string{"ls", "-la", "non-existent-file"},
				env: nil,
			},
			wantReturnCode: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotReturnCode := RunCmd(tt.args.cmd, tt.args.env); gotReturnCode != tt.wantReturnCode {
				t.Errorf("RunCmd() = %v, want %v", gotReturnCode, tt.wantReturnCode)
			}
		})
	}
}
