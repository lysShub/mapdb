package main

import (
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkComprehensive(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		if err = Comprehensive(strconv.Itoa(i)); err != nil {
			fmt.Println("error", err)
		}
	}
}
