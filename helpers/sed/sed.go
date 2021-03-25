package sed

import (
	"errors"
	"strconv"
	"strings"
)

type SedStruct struct {
	original   string
	replace    string
	repetation int
}

func ParseSed(input string) (*SedStruct, error) {

	messageWithoutCommand := strings.Replace(input, "/sed", "", 1)
	messageWithoutCommand = strings.Trim(messageWithoutCommand, " ") //Removing leading and trailing whitespaces

	// /sed ladka/ladki/1

	splitInput := strings.Split(messageWithoutCommand, "/")
	var sedstruct *SedStruct
	if len(splitInput) == 3 {
		original := splitInput[0]
		replace := splitInput[1]
		repeat, err := strconv.Atoi(splitInput[2])
		if err != nil {
			return nil, err
		}
		sedstruct = &SedStruct{original: original, replace: replace, repetation: repeat}

	} else if len(splitInput) == 2 {
		original := splitInput[0]
		replace := splitInput[1]
		sedstruct = &SedStruct{original: original, replace: replace, repetation: 1}
	} else {
		return nil, errors.New("Invalid string")
	}

	return sedstruct, nil

}

func Sed(repliedmessage string, input string) (string, error) {

	if sedstruct, err := ParseSed(input); err != nil {
		return "", err
	} else {
		newmessage := strings.Replace(repliedmessage, sedstruct.original, "*"+sedstruct.replace+"*", sedstruct.repetation)
		return newmessage, nil
	}
}
