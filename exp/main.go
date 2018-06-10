package main

import (
	"fmt"

	"github.com/jlindauer/usegolang/hash"
)

func main() {
	hmac := hash.NewHMAC("my-secret-key")
	fmt.Println(hmac.Hash("this is my string to hash"))
}
