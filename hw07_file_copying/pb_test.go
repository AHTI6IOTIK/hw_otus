package main

import (
	"bytes"
	"github.com/AHTI6IOTIK/hw_otus/hw07_file_copying/progressbar"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestNewProgressBar(t *testing.T) {
	type prepare func(t *testing.T, pb *progressbar.ProgressBar)
	type args struct {
		width int
		out   io.ReadWriter
		size  int
	}
	tests := []struct {
		name      string
		args      args
		wantOut   string
		prepare   prepare
		wantErr   error
		isWantErr bool
	}{
		{
			name:      "check invalid width with 0",
			args:      args{width: 0, out: &bytes.Buffer{}},
			wantErr:   progressbar.ErrInvalidPbWidth,
			isWantErr: true,
		},
		{
			name:      "check invalid width with -10",
			args:      args{width: -10, out: &bytes.Buffer{}},
			wantErr:   progressbar.ErrInvalidPbWidth,
			isWantErr: true,
		},
		{
			name:      "check invalid out",
			args:      args{width: 10, out: nil},
			wantErr:   progressbar.ErrInvalidPbOut,
			isWantErr: true,
		},
		{
			name:    "check success",
			args:    args{width: 10, out: &bytes.Buffer{}, size: 50},
			wantOut: "\r[#         ] 10%",
			prepare: func(t *testing.T, pb *progressbar.ProgressBar) {
				_, err := pb.Write([]byte{1, 2, 3, 4, 5})
				if err != nil {
					t.Error(err)
					return
				}
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := progressbar.NewProgressBar(tt.args.width, tt.args.out)
			if tt.isWantErr {
				require.Equal(t, tt.wantErr, err)
				return
			}

			got.Size(tt.args.size)
			tt.prepare(t, got)
			r := make([]byte, tt.args.size+10)
			n, err := tt.args.out.Read(r)
			if err != nil {
				t.Error(err)
				return
			}
			require.Equal(t, []byte(tt.wantOut), r[:n])
		})
	}
}
