package utils

import "testing"

func TestMatchLabels(t *testing.T) {
	type args struct {
		sourceLabels map[string]string
		targetLabels map[string]string
	}
	labelsOne := map[string]string{
		"component": "my-component",
	}
	labelsTwo := map[string]string{
		"component": "my-component2",
	}
	labelsDifferent := map[string]string{
		"app": "my-component3",
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "source and target labels match",
			args: args{
				sourceLabels: labelsOne,
				targetLabels: labelsOne,
			},
			want: true,
		},
		{
			name: "source and target labels don't match",
			args: args{
				sourceLabels: labelsOne,
				targetLabels: labelsTwo,
			},
			want: false,
		},
		{
			name: "source and target labels are different",
			args: args{
				sourceLabels: labelsOne,
				targetLabels: labelsDifferent,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchLabels(tt.args.sourceLabels, tt.args.targetLabels); got != tt.want {
				t.Errorf("MatchLabels() = %v, want %v", got, tt.want)
			}
		})
	}
}
