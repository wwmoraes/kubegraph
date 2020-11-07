package utils

import "testing"

func TestContainsString(t *testing.T) {
	type args struct {
		array   []string
		element string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "returns false when array doesn't contain element",
			args: args{
				array:   []string{"a", "b", "c"},
				element: "x",
			},
			want: false,
		},
		{
			name: "returns true when array contains element",
			args: args{
				array:   []string{"a", "b", "c"},
				element: "a",
			},
			want: true,
		},
		{
			name: "returns false on empty array",
			args: args{
				array:   []string{},
				element: "x",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsString(tt.args.array, tt.args.element); got != tt.want {
				t.Errorf("ContainsString() = %v, want %v", got, tt.want)
			}
		})
	}
}
