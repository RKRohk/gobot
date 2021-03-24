package reminders_test

import (
	"testing"

	"github.com/rkrohk/gobot/helpers/reminders"
	"github.com/stretchr/testify/assert"
)

func TestRegex(t *testing.T) {

	tests := []string{"01/04/2021 2:51PM Happy Birthday", "01/04/2021 2:51PM Hi", "01/04/2021 12:51PM Exam today", "25/02/2021 12:51PM Assignment Due"}
	expected := []string{"01/04/2021 2:51PM", "01/04/2021 2:51PM", "01/04/2021 12:51PM", "25/02/2021 12:51PM"}

	for i, test := range tests {
		indices := reminders.GetDateIndices(test)
		a, b := indices[0], indices[1]
		assert.Equal(t, test[a:b], expected[i])
	}

}
