package test

import "testing"

func IsPalindrome(letters string) bool {
	list := make(map[byte]bool, 10000)
	length := len(letters)
	n := length / 2
	for i := 0; i < n; i++ {
		if letters[i] != letters[length-1-i] {
			list[letters[i]] = true
			return false
		}
	}
	return true
}

func BenchmarkIsPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPalindrome("A man, a plan, a canal: Panama")
	}
}
