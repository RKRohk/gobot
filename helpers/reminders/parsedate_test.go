package reminders_test

import (
	"testing"
	"time"

	"github.com/rkrohk/gobot/helpers/reminders"
	"github.com/stretchr/testify/assert"
)

func TestParseDate(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Kolkata")

	dates := []string{"01/04/2021 2:51PM IST", "1/04/2021 02:51PM IST", "01/4/2021 02:51PM IST", "1/4/2021 02:51PM IST"}

	actualDate := time.Date(2021, time.April, 01, 14, 51, 0, 0, loc)

	for _, date := range dates {
		parsedDate, err := reminders.ParseDate(date)
		if err != nil {
			t.Error("Error in parsing date", err)
		} else {
			assert.Equal(t, parsedDate, actualDate)
		}
	}

	dates = []string{"1/4/2021 12:51PM IST", "01/04/2021 12:51PM IST"}
	actualDate = time.Date(2021, time.April, 01, 12, 51, 0, 0, loc)

	for _, date := range dates {
		parsedDate, err := reminders.ParseDate(date)
		if err != nil {
			t.Error("Error in parsing date", err)
		} else {
			assert.Equal(t, parsedDate, actualDate)
		}
	}

	parsedDate, _ := reminders.ParseDate("29/03/2021 6:56PM IST")
	assert.Equal(t, parsedDate, time.Date(2021, 3, 29, 18, 56, 0, 0, loc))

}
