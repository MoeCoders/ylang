package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
)

type (
	codeList []string
	code     []string
)

// DeleteSlice 删除指定元素。
func (s code) deleteSlice(elem string) code {
	j := 0
	for _, v := range s {
		if v != elem {
			s[j] = v
			j++
		}
	}
	return s[:j]
}

func (s *codeList) readFile(filePath string) error {
	f, err := os.Open(filePath)
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err) // 或设置到函数返回值中
		}
	}()
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if len(line) > 0 && line != "\n" {
			*s = append(*s, line)
		}
		if err == io.EOF {
			return nil
		}
	}
}

// 词法分析器(lexer)：将代码字符串解析为token切片
func lex(input string) (code, error) {
	pattern := `\d+(\.\d+)?|[+\-*/\^()]`
	re := regexp.MustCompile(pattern)
	tokens := re.FindAllString(input, -1)
	matchRe, err := regexp.Compile(`^[+\\-]$`)
	if err != nil {
		return nil, err
	}
	nRe, err := regexp.Compile(`^[()+\\-\\*/\\^]$`)
	if err != nil {
		return nil, err
	}
	for i, v := range tokens {
		match := matchRe.MatchString(v)
		if err != nil {
			return nil, err
		}
		if i == 0 && match {
			tokens[1] = tokens[0] + tokens[1]
			tokens = tokens[1:]
			continue
		} else if i > 0 && match {
			n := nRe.MatchString(tokens[i-1])
			if err != nil {
				return nil, err
			}
			if n {
				tokens[i+1] = tokens[i] + tokens[i+1]
				tokens = append(tokens[:i], tokens[i+1:]...)
			}
		}
	}
	return tokens, nil
}

// 解释器(evaluator)：读取 token 并进行计算
func eval(tokens code) (float64, error) {
	tokens = tokens.deleteSlice("")
	numStack := []float64{}
	opStack := []string{}
	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/", "^":
			for len(opStack) > 0 && isHigherOrEqualPrecedence(opStack[len(opStack)-1], token) {
				if err := evaluateTopOfStack(&numStack, &opStack); err != nil {
					return 0, err
				}
			}
			opStack = append(opStack, token)
		case "(":
			opStack = append(opStack, token)
		case ")":
			for len(opStack) > 0 && opStack[len(opStack)-1] != "(" {
				if err := evaluateTopOfStack(&numStack, &opStack); err != nil {
					return 0, err
				}
			}
			if len(opStack) == 0 {
				return 0, fmt.Errorf("mismatched parentheses")
			}
			opStack = opStack[:len(opStack)-1]
		default:
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			numStack = append(numStack, num)
		}
	}
	for len(opStack) > 0 {
		if opStack[len(opStack)-1] == "(" {
			return 0, fmt.Errorf("mismatched parentheses")
		}
		if err := evaluateTopOfStack(&numStack, &opStack); err != nil {
			return 0, err
		}
	}
	if len(numStack) != 1 || len(opStack) > 0 {
		return 0, fmt.Errorf("invalid expression")
	}
	return numStack[0], nil
}

func isHigherOrEqualPrecedence(op1, op2 string) bool {
	if op2 == "^" {
		return false
	}
	if op1 == "*" || op1 == "/" {
		return true
	}
	if op2 == "*" || op2 == "/" {
		return false
	}
	if op1 == "(" {
		return false
	}
	if op2 == "(" {
		return true
	}
	return true
}

func evaluateTopOfStack(numStack *[]float64, opStack *[]string) error {
	if len(*opStack) < 1 {
		return fmt.Errorf("invalid expression")
	}
	op := (*opStack)[len(*opStack)-1]
	*opStack = (*opStack)[:len(*opStack)-1]
	nStack := *numStack
	num2 := nStack[len(nStack)-1]
	nStack = nStack[:len(nStack)-1]
	num1 := nStack[len(nStack)-1]
	nStack = nStack[:len(nStack)-1]
	var result float64
	switch op {
	case "^":
		result = math.Pow(num1, num2)
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*":
		result = num1 * num2
	case "/":
		if num2 == 0 {
			return fmt.Errorf("divide by zero")
		}
		result = num1 / num2
	}
	*numStack = append(nStack, result)
	return nil
}

func (s *codeList) run() {
	for _, v := range *s {
		tokens, err := lex(v)
		if err != nil {
			fmt.Println(err)
			return
		}
		result, err := eval(tokens)
		if err != nil {
			fmt.Println("error: ", err)
			return
		}
		fmt.Printf("%f\n", result)
	}
}

func cli() {
	var codes codeList = make(codeList, 1)
	for {
		fmt.Printf(">>> ")
		var c string
		fmt.Scanln(&c)
		if c == "exit" {
			fmt.Println("Goodbye!")
			return
		}
		codes[0] = c
		codes.run()
	}
}

func help() {
	fmt.Printf("")
}

func main() {
	if len(os.Args) > 2 {
		switch os.Args[1] {
		case "run":
			var codes codeList
			err := codes.readFile(os.Args[2])
			if err != nil {
				fmt.Println(err)
				return
			}
			codes.run()
			return
		default:
			help()
		} }else if len(os.Args) == 1{
            cli()
		}
}
