package main

import "huseynovvusal/gojudge/internal/judge"

func main() {
	// result, err := judge.RunCode("python", "n = 1\nwhile True:\n\tn += 1")
	result, err := judge.RunCode("python", "n = 1\nwhile (n < 100):\n\tn += 1; print(n)")
	if err != nil {
		panic(err)
	}

	println("RESULT:" + result)

}
