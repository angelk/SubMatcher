package main

import "testing"

func TestMatchScore(t *testing.T) {
	type matchTestCase struct {
		A string
		B string
		LenMatch int
		Score int
	}

	testCases := make([]matchTestCase, 0)
	testCases = append(
		testCases,
		matchTestCase{
			A: "aaabbbccc",
			B: "vvvfffccc",
			LenMatch: 3,
			Score: 1,
		},
		matchTestCase{
			A: "aaabbbccc",
			B: "qqqwwweee",
			LenMatch: 3,
			Score: 0,
		},
		matchTestCase{
			A: "",
			B: "qqqwwweee",
			LenMatch: 3,
			Score: 0,
		},
		matchTestCase{
			A: "aaawwwaaawww",
			B: "aaww",
			LenMatch: 3,
			Score: 4,
		},
	)

	for _, testCase := range testCases {
		score := getMatchScore(testCase.A, testCase.B, testCase.LenMatch)
		if score != testCase.Score {
			t.Error("Expecting score of 3, got ", score)
		}
	}
}
