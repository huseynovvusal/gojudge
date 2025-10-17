package main

import "huseynovvusal/gojudge/internal/judge"

func main() {
	// result, err := judge.RunCode("python", "n = 1\nwhile True:\n\tn += 1")
	// result, err := judge.RunCode("python", "n = 1\nwhile (n < 100):\n\tn += 1; print(n)")
	result, err := judge.RunCode("cpp", "#include <iostream>\nint main() {\n    int n = 1;\n    while (n < 5000) {\n        n++;\n        std::cout << n << std::endl;\n    }\n    return 0;\n}")
	if err != nil {
		panic(err)
	}

	println("RESULT:" + result)

}
