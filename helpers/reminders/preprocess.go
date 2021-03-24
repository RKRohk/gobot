package reminders

import "regexp"

func GetDateIndices(text string) []int {
	dateRegex := regexp.MustCompile(`[012]?\d\/0?\d\/20\d\d [0|1]?\d:\d\d[A?P]M`)
	return dateRegex.FindStringIndex(text)
}
