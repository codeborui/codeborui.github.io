package main

import (
	"fmt"
)

func main() {
	nums := []int{1, 2, 2, 2, 3, 4}
	fmt.Println(searchRange(nums, 2))
}

func searchRange(nums []int, target int) []int {
	if nums == nil || len(nums) == 0 {
		return []int{-1, -1}
	}
	return []int{searchFirst(nums, target), searchLast(nums, target)}
}

func searchFirst(nums []int, target int) int {
	left, right, mid := 0, len(nums)-1, 0
	for left <= right {
		mid = left + (right-left)>>1
		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	if left < len(nums) && nums[left] == target {
		return left
	}
	return -1
}

func searchLast(nums []int, target int) int {
	left, right, mid := 0, len(nums)-1, 0
	for left <= right {
		mid = left + (right-left)>>1
		if nums[mid] > target {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	if right >= 0 && nums[right] == target {
		return right
	}
	return -1
}
