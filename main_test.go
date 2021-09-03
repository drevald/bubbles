package main

import (
	"testing"
	"fmt"
)

func TestRotateMatrix (t *testing.T) {
	m := Matrix{cells:[]int{1,1,1,0,1,0}, width:3, height:2}
	fmt.Println(m)
	fmt.Println("Rotating right")
	m.RotateRight()
	fmt.Println(m)
	fmt.Println("Rotating left")
	m.RotateLeft()
	fmt.Println(m)
	fmt.Println("Testing")
}