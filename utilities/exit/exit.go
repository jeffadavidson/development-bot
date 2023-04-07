package exit

import (
	"fmt"
	"os"
)

func ExitSuccess() {
	os.Exit(0)
}

func ExitError(err error) {
	fmt.Println(err.Error())
	os.Exit(-1)
}
