package sed_test

import (
	"testing"

	"github.com/rkrohk/gobot/helpers/sed"
)

func TestSedParser(t *testing.T) {
	inputStrings := []string{"/sed man/woman", "/sed man/woman/2", "/sed smart/stupid/5", "/sed this guy is very smart/this guy is not very smart"}

	outputs := make([]*sed.SedStruct, len(inputStrings))
	for i, input := range inputStrings {
		output, err := sed.ParseSed(input)
		if err != nil {
			t.Error("Error parsing sed")
		}
		outputs[i] = output
	}
}

func TestSed(t *testing.T) {
	op, err := sed.Sed("good morning children", "/sed children/child")
	if err == nil {
		t.Log(op)
	}
}
