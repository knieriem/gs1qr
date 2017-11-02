package gs1qr

import (
	"fmt"
	"strings"

	"github.com/knieriem/gs1qr/ai"
)

func Example() {
	list, err := ai.ParseSeq("(01)03453120000011(8200)http://abc.net(10)XYZ(410)9501101020917")
	if err != nil {
		fmt.Println(err)
		return
	}
	el := ConvertElements(list)
	fmt.Println(strings.Join(el.Strings(), ""))
	// Output: <FNC1>01034531200000118200http://abc.net<GS>10XYZ%4109501101020917
	return
}
