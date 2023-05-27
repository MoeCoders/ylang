package main

import (
	"fmt"
	"regexp"
)

func main() {
	equation := "1.5+2*(3-4)/-0.25"
	pattern := `\d+(\.\d+)?|[+\-*/\^()]`
	re := regexp.MustCompile(pattern)
	tokens := re.FindAllString(equation, -1)
	fmt.Println(tokens)
	for i, v := range tokens {
		match, err := regexp.MatchString("^[+\\-]$", v)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
        if i == 0 && match {
            tokens[1] = tokens[0] + tokens[1]
            tokens = tokens[1:]
            continue
        }
		n, err := regexp.MatchString("^[()+\\-\\*/\\^]", tokens[i-1])
        if err != nil {
            fmt.Println(err)
        }
        if match && n {
            
        }
	}
}
