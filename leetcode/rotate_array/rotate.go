package rotate_array

import "math/rand"

func rotateCopy(nums []int, k int) {
	k %= len(nums)

	copy(nums, append(nums[len(nums)-k:], nums[:len(nums)-k]...))
}

func rotateInPlace(nums []int, k int) {
	k %= len(nums)
	reverse(nums)
	reverse(nums[:k])
	reverse(nums[k:])
}

func rotateGnGn(nums []int, k int) {
	lenNums := len(nums)
	kSanitized := k % lenNums

	result := append(nums[lenNums-kSanitized:], nums[:lenNums-kSanitized]...)
	for i := 0; i < lenNums; i++ {
		nums[i] = result[i]
	}
}

func reverse(nums []int) {
	i := 0
	j := len(nums) - 1
	for i < j {
		nums[i], nums[j] = nums[j], nums[i]
		i++
		j--
	}
}

func random(min int, max int) int {
	return min + rand.Intn(int(max-min))
}

func RandomArray(size int, min int, max int) []int {
	var array = make([]int, size)
	for i := 0; i < size; i++ {
		array[i] = random(min, max)
	}
	return array
}
