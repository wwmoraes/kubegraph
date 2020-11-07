package utils

import (
	"reflect"
	"testing"
)

func TestStringChunks(t *testing.T) {
	type args struct {
		str       string
		chunkSize int
	}

	testString := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	testChunks := []string{
		"ABCDEFGHIJKLM",
		"NOPQRSTUVWXYZ",
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "string lower then chunkSize",
			args: args{
				str:       testString,
				chunkSize: len(testString),
			},
			want: []string{testString},
		},
		{
			name: "string has been chunked",
			args: args{
				str:       testString,
				chunkSize: len(testString) / 2,
			},
			want: testChunks,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringChunks(tt.args.str, tt.args.chunkSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringChunks() = %v, want %v", got, tt.want)
			}
		})
	}
}
