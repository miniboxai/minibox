package pretty

import (
	"fmt"

	p "github.com/kr/pretty"
)

func Formatter(x interface{}) (f fmt.Formatter) {
	return p.Formatter(x)
}
