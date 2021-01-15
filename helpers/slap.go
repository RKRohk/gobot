package helpers

import (
	"fmt"
)

//Slap takes two users as argument and slaps one
func Slap(user1 string, user2 string) string {
	pick, err := GetSlapStrings()
	if err != nil {
		return "database error occurred"
	}
	return fmt.Sprintf(pick, user1, user2)
}
