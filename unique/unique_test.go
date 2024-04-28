package unique

import (
	"fmt"
	"testing"
)

func TestUniqueValueUUID(t *testing.T) {
	for i := 0; i < 10; i++ {
		result := UniqueValueUUID()
		fmt.Printf("%v\r\n", result)
	}
	for i := 0; i < 10; i++ {
		result := UniqueValueUUIDNew()
		fmt.Printf("%v\r\n", result)
	}
}
