package pkg

import (
	pkg2 "cs/testvar/pkg2"
	"fmt"
)

func call(sa pkg2.ISA) {
	sa.Set(3)
	fmt.Println(sa.Get())
}
