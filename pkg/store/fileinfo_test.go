package store

import (
	"testing"
)

func Test_fileinfo(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "Test_fileinfo",
			args: []string{
				// "MAVRYK_HANGZHOUNET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full",
				"MAVRYK_BASENET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.rolling",
				"MAVRYK_MAINNET-BL4p3YRfxhiQP16PsuzFBbph8QqVcNN3qu42r5JgNgdaw3xW81g-2396664.rolling",
				"MAVRYK_MAINNET-BMBDsvNoA4wr4VANmUJfMPPEpCKKBYY7xBoYfhJkUuoGk54GYPa-4593763.rolling",
			},
			want: []string{
				// "hangzhounet",
				"basenet",
				"mainnet",
				"mainnet",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, want := range tt.want {
				got := getInfoFromfilename(tt.args[i])
				if got.ChainName != want {
					t.Errorf("got: %v, want: %v", got, want)
				}
			}
		})
	}
}
