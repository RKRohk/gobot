package reminders_test

import (
	"testing"
	"time"

	"github.com/rkrohk/gobot/helpers/reminders"
	"github.com/stretchr/testify/assert"
)

func TestParseDate(t *testing.T) {

	dates := []string{"01/04/2021 2:51PM", "1/04/2021 02:51PM", "01/4/2021 02:51PM", "1/4/2021 02:51PM"}

	actualDate := time.Date(2021, time.April, 01, 14, 51, 0, 0, time.UTC)

	for _, date := range dates {
		parsedDate, err := reminders.ParseDate(date)
		if err != nil {
			t.Error("Error in parsing date", err)
		} else {
			assert.Equal(t, parsedDate, actualDate)
		}
	}

	dates = []string{"1/4/2021 12:51PM", "01/04/2021 12:51PM"}
	actualDate = time.Date(2021, time.April, 01, 12, 51, 0, 0, time.UTC)

	for _, date := range dates {
		parsedDate, err := reminders.ParseDate(date)
		if err != nil {
			t.Error("Error in parsing date", err)
		} else {
			assert.Equal(t, parsedDate, actualDate)
		}
	}

}
