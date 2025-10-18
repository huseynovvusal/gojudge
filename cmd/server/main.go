package main

import (
	"fmt"
	"huseynovvusal/gojudge/internal/judge"
)

func main() {

	result, err := judge.RunCode("python", "print('Hello, World!')", 2, 256, 1)
	// result, err := judge.RunCode("cpp", "#include <iostream>\nint main() { std::cout << \"Hello, C++ World!\" << std::endl; return 0; }", 2, 256, 1)
	if err != nil {
		fmt.Printf("Error: %v\nOutput: %s\n", err, result.Output)
		return

	}

	println("Output:\n" + result.Output)
	println("Execution Time (ms):", result.ExecutionMs)

}
