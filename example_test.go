package indexnow_test

import (
	"fmt"

	"github.com/aquilax/indexnow"
)

func ExampleGetSingleSubmitUrl() {
	var u string
	u = indexnow.GetSingleSubmitUrl("example.com", "aabbccddeeff", "", "https://www.example.com/")
	fmt.Println(u)
	u = indexnow.GetSingleSubmitUrl("example.com", "", "https://www.example.com/key.txt", "https://www.example.com/")
	fmt.Println(u)
	// Output:
	// https://example.com/indexnow?key=aabbccddeeff&url=https%3A%2F%2Fwww.example.com%2F
	// https://example.com/indexnow?key=&keyLocation=https%3A%2F%2Fwww.example.com%2Fkey.txt&url=https%3A%2F%2Fwww.example.com%2F
}
