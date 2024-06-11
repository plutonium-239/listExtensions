package scrollinglist

// Find an index such that 0 to index just (over)fills valueToReach
// i.e. the appropriate value - 1
func PrefixSumBreak(list []int, valueToReach int) (int, int) {
	sum := 0
	for i := 0; i < len(list); i++ {
		sum += list[i]
		if sum >= valueToReach {
			return i, sum - valueToReach
		}
	}
	return len(list), 0
}

// Find an index such that index to len(list) just (over)fills valueToReach
// i.e. the appropriate value + 1
func SuffixSumBreak(list []int, valueToReach int) (int, int) {
	sum := 0
	for i := 1; i < len(list); i++ {
		sum += list[len(list)-i]
		if sum >= valueToReach {
			return i, sum - valueToReach
		}
	}
	return len(list), 0
}
