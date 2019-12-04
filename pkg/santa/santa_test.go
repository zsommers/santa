package santa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPersonToString(t *testing.T) {
	// For making references
	jane := "Jane"
	fred := "Fred"

	cases := []struct {
		name     string
		input    Person
		expected string
	}{
		{
			"Basic",
			Person{
				Name:      "John",
				Mobile:    "+15551234567",
				Excluded:  []string{"Joe", "Fred"},
				Recipient: &jane,
				Santa:     &fred,
			},
			"John\tRecipient: Jane\tExcluded: [Joe Fred]",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.input.String())
		})
	}
}
