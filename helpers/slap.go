package helpers

import (
	"fmt"
	"math/rand"
)

//Slap takes two users as argument and slaps one
func Slap(user1 string, user2 string) string {
	arr := []string{"%s burned %s to death", "%s hit %s with a chair", "%s hit %s around a bit with a large trout", "%s burnt %s to a crisp in backyard bbq"}
	randomIndex := rand.Intn(len(arr))
	pick := arr[randomIndex]
	return fmt.Sprintf(pick, user1, user2)
}
