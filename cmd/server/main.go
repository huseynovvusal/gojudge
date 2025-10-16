package main

import "huseynovvusal/gojudge/internal/judge"

func main() {
	result, err := judge.RunCode("python", "print('Hello, World!')")
	if err != nil {
		panic(err)
	}

	println("RESULT:" + result)

}
