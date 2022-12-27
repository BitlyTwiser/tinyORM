package tinyorm_test

import "testing"

func TestInitializeDatabase(t *testing.T) {
	tests := []struct {
		name string
		have string
		want string
	}{
		{
			name: "",
			have: "",
			want: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.have != test.want {
				t.Errorf("test %s failed. Have: %v Want: %v", test.name, test.have, test.want)
			}
		})
	}
}
