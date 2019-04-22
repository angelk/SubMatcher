package main

func getMatchScore(a string, b string, diffLen int) int {
	matchScore := 0

	aLen := len(a)
	bLen := len(b)
	for i := 0; i+diffLen-1 < aLen; i++ {
		aPart := a[i : i+diffLen]
		for j := 0; j+diffLen-1 < bLen; j++ {
			bPart := b[j : j+diffLen]
			if aPart == bPart {
				matchScore++
			}
		}
	}

	return matchScore
}
