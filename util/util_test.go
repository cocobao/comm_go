package util

import (
	"fmt"
	"testing"
	"time"
)

func TestGetDate(t *testing.T) {
	fmt.Println(time.Now().Format("01-02"))
}
