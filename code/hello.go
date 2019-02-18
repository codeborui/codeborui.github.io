package main

import (
	"fmt"
)

func main() {
	nums := []int{1}
	fmt.Println(combinationSum(nums, 1))
}

func combinationSum(candidates []int, target int) [][]int {
    return [][]int{{0}}
}
