package main

import (
	"fmt"
)

func main() {
	for _, x := range servicesErrMap {
		fmt.Printf("%d - %s - %s\n", x.Code, x.Body.ErrorType, x.Body.ErrorMessage)
	}
	for _, x := range domainErrMap {
		fmt.Printf("%d - %s - %s\n", x.Code, x.Body.ErrorType, x.Body.ErrorMessage)
	}
}
