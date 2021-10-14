package functioncaller

import (
	"fmt"
	"runtime"
)

func PrintFuncName() string {
	fpcs := make([]uintptr, 1)

	// Skip 2 levels to get the caller
	n := runtime.Callers(2, fpcs)
	if n == 0 {
		fmt.Println("MSG: NO CALLER")
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		fmt.Println("MSG CALLER WAS NIL")
	}

	if caller == nil {
		return ""
	}

	return caller.Name()
}
