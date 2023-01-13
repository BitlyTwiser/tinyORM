package sqlbuilder

import "testing"

func TestSnakeCase(t *testing.T) {
	tests := map[string]struct {
		have string
		want string
	}{
		"Test with ID enhanced": {have: "ID", want: "id"},
		"Test with ID":          {have: "Id", want: "id"},
		"Simple Test":           {have: "Testing", want: "testing"},
		"Secondary Test":        {have: "HowNowBrownCow", want: "how_now_brown_cow"},
		"Advanced Test":         {have: "TableWithNoID", want: "table_with_no_id"},
		"Diff casing Test":      {have: "imma_Table", want: "imma__table"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			lowerSnakeCaseString := lowerSnakeCase(test.have)
			if test.want != lowerSnakeCaseString {
				t.Fatalf("Wanted: %s - Have: %s", test.want, lowerSnakeCaseString)
			}
		})
	}
}
