package reminders

import "regexp"

//GetDateIndices finds the start and end indices of the date and time from a given input message
func GetDateIndices(text string) []int {
	dateRegex := regexp.MustCompile(`[012]?\d\/0?\d\/20\d\d [0|1]?\d:\d\d[A?P]M [A-Z]{3}`)
	return dateRegex.FindStringIndex(text)
}
