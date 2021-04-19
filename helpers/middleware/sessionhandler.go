package middleware

/*The aim of this entire thing is to check which command
the user previously used, e.g makepdf. So, the function will
call makepdf command for the user until the user gives a /done command
*/

var imageArr = make(map[UserChatID][]interface{})

func MakePdf(s *UserSession) {
	if images := imageArr[s.ChatID]; images != nil {
		imageArr[s.ChatID] = append(images, 1)
	} else {
		imageArr[s.ChatID] = []interface{}{1, 2, 3}
	}
	print(imageArr)
}
