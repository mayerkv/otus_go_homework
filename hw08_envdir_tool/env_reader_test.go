package main

import (
	"os"
	"reflect"
	"testing"
)

const (
	chmod0    = 0b0
	chmod0664 = 0b000_110_110_100
)

func setup(t *testing.T) {
	t.Helper()
	if err := os.Chmod("testdata/envnopermission/NO_PERM", chmod0); err != nil {
		t.Error(err)
	}
}

func teardown(t *testing.T) {
	t.Helper()
	if err := os.Chmod("testdata/envnopermission/NO_PERM", chmod0664); err != nil {
		t.Error(err)
	}
}

func TestReadDir(t *testing.T) {
	setup(t)
	defer teardown(t)

	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    Environment
		wantErr bool
	}{
		{
			name: "read testdata/env",
			args: args{"testdata/env"},
			want: Environment{
				"BAR":   {"bar", false},
				"EMPTY": {"", false},
				"FOO":   {"   foo\nwith new line", false},
				"HELLO": {"\"hello\"", false},
				"UNSET": {"", true},
			},
			wantErr: false,
		},
		{
			name:    "read non-existent file",
			args:    args{"not_existent_file"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "skip directory in env path",
			args: args{"testdata/envwithdir"},
			want: Environment{
				"BAR":   {"bar", false},
				"EMPTY": {"", false},
				"FOO":   {"   foo\nwith new line", false},
				"HELLO": {"\"hello\"", false},
				"UNSET": {"", true},
			},
			wantErr: false,
		},
		{
			name:    "read with invalid symbol in filename",
			args:    args{"testdata/envwithinvalidsymbol"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "read when no permission to read the file",
			args:    args{"testdata/envnopermission"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDir(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}
