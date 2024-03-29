package libwekan

import (
	"math/rand"
	"time"
)

func newId() string {
	return newIdN(17)
}

func newId6() string {
	return newIdN(6)
}

func newIdN(n int) string {
	chars := "123456789ABCDEFGHJKLMNPQRSTWXYZabcdefghijkmnopqrstuvwxyz"
	l := len(chars)
	var digits []byte
	for i := 0; i < n; i++ {
		digit := rand.Intn(l)
		digits = append(digits, chars[digit])
	}
	return string(digits)
}

func uniq[Element comparable](array []Element) []Element {
	m := make(map[Element]struct{})
	for _, element := range array {
		m[element] = struct{}{}
	}
	var set []Element
	for element := range m {
		set = append(set, element)
	}
	return set
}

func intersect[E comparable](elementsA []E, elementsB []E) (both []E, onlyA []E, onlyB []E) {
	for _, elementA := range elementsA {
		foundBoth := false
		for _, elementB := range elementsB {
			if elementA == elementB {
				both = append(both, elementA)
				foundBoth = true
			}
		}
		if !foundBoth {
			onlyA = append(onlyA, elementA)
		}
	}

	for _, elementB := range elementsB {
		foundBoth := false
		for _, elementA := range elementsA {
			if elementA == elementB {
				foundBoth = true
			}
		}
		if !foundBoth {
			onlyB = append(onlyB, elementB)
		}
	}
	return both, onlyA, onlyB
}

func mapSlice[T any, M any](a []T, f func(T) M) []M {
	n := make([]M, len(a))
	for i, e := range a {
		n[i] = f(e)
	}
	return n
}

func toMongoTime(t time.Time) time.Time {
	return t.In(time.UTC).Truncate(time.Millisecond)
}

func contains[Element comparable](elements []Element, element Element) bool {
	for _, actual := range elements {
		if element == actual {
			return true
		}
	}
	return false
}

func selectSlice[Element any](slice []Element, filter func(Element) bool) []Element {
	var accepted []Element
	for _, element := range slice {
		if filter(element) {
			accepted = append(accepted, element)
		}
	}
	return accepted
}
