package main

import (
	"fmt"
	"strconv"
	"strings"
)

const IsLeftParenthesis = 1
const IsRightParenthesis = 2
const IsNotParenthesis = 0
const IsMultiplication = 10
const IsDivision = 20
const IsAddition = 30
const IsSubtraction = 40
const IsNotOperation = 0
const IsPoint = 100
const IsNotSeparator = 0

func IsNumber(char byte) bool {
	const numbers = "1234567890"
	for index, _ := range numbers {
		if numbers[index] == char {
			return true
		}
	}
	return false
}

func IsParenthesis(char byte) int {
	if string(char) == "(" {
		return IsLeftParenthesis
	}
	if string(char) == ")" {
		return IsRightParenthesis
	}
	return IsNotParenthesis
}

func IsOperator(char byte) int {
	if string(char) == "*" {
		return IsMultiplication
	} else if string(char) == "/" {
		return IsDivision
	} else if string(char) == "+" {
		return IsAddition
	} else if string(char) == "-" {
		return IsSubtraction
	}
	return IsNotOperation
}

func MathOp(sliceofnums []float64, operation int) (float64, error) {
	var result float64
	var lengthofslice int = len(sliceofnums)

	num1 := sliceofnums[lengthofslice-2]
	num2 := sliceofnums[lengthofslice-1]

	switch {
	case operation == IsAddition:
		result = num1 + num2
	case operation == IsSubtraction:
		result = num1 - num2
	case operation == IsMultiplication:
		result = num1 * num2
	case operation == IsDivision:
		if num2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		result = num1 / num2
	}
	return result, nil
}

func IsSeparator(char byte) int {
	if string(char) == "." {
		return IsPoint
	}
	return IsNotSeparator
}

func GetPryority(operator int) int {
	mapofoperators := map[int]int{
		IsMultiplication: 2,
		IsDivision:       2,
		IsAddition:       1,
		IsSubtraction:    1,
	}
	pryority := mapofoperators[operator]
	return pryority
}

func ExtractNum(Expression string, indexofnum int) (string, int) {
	var num string
	var index int
	var length int = len(Expression)
	for nextnotnumindex := indexofnum; nextnotnumindex < length; nextnotnumindex++ {
		if IsNumber(Expression[nextnotnumindex]) || IsSeparator(Expression[nextnotnumindex]) != 0 {
			num += string(Expression[nextnotnumindex])
		}
		if !IsNumber(Expression[nextnotnumindex]) && IsSeparator(Expression[nextnotnumindex]) == 0 {
			return num, nextnotnumindex
		}
		index = nextnotnumindex
	}
	return num, index
}

func PopNum(sliceofnums []float64, indexofnum int) ([]float64, error) {

	var newsliceofnums []float64

	if indexofnum > len(sliceofnums) {
		return sliceofnums, fmt.Errorf("index of num > length of slice of nums, %d", indexofnum)
	}
	if indexofnum < 0 {
		return sliceofnums, fmt.Errorf("index of num < 0, %d", indexofnum)
	}

	newsliceofnums = append(sliceofnums[:indexofnum], sliceofnums[indexofnum+1:]...)

	return newsliceofnums, nil
}

func PopOp(opslice []int, opindex int) ([]int, error) {

	var newopslice []int

	if opindex > len(opslice) {
		return opslice, fmt.Errorf("index of operator > length of slice of operator, %d", opindex)
	}
	if opindex < 0 {
		return opslice, fmt.Errorf("index of operator < 0, %d", opindex)
	}

	newopslice = append(opslice[:opindex], opslice[opindex+1:]...)

	return newopslice, nil
}

func IsCorrectExpression(Expression string) (bool, error) { //Проверка на правильность заданной строки
	Expression = strings.Replace(Expression, " ", "", -1)
	if Expression == "" {
		return false, fmt.Errorf("empty expression")
	}
	correctexpression := true
	expressionlength := len(Expression)
	countleftparenthesis := 0
	countrightparenthesis := 0
	for index, _ := range Expression {
		if index < expressionlength-1 {
			switch {
			case !IsNumber(Expression[index]) && IsParenthesis(Expression[index]) == 0 && IsOperator(Expression[index]) == 0 && IsSeparator(Expression[index]) == 0: //Недопустимые символы
				correctexpression = false
				//fmt.Printf("incorrect symbol, char %d. Allowed only: %s ", index, "1234567890.*/+-()")
			case index == 0 && !IsNumber(Expression[index]) && IsParenthesis(Expression[index]) == 0: //Запрещенная последовательность "выражение начинается не числом и не скобкой"
				correctexpression = false
				//fmt.Printf(`wrong sequence "non-number character": char %d `, index)
			case IsOperator(Expression[index]) != 0 && IsOperator(Expression[index+1]) != 0: //Запрещенная последовательность "оператор->оператор"
				correctexpression = false
				//fmt.Printf(`wrong sequence "operation sign->operation sign": chars %d, %d `, index, index+1)
			case IsSeparator(Expression[index]) != 0 && IsSeparator(Expression[index+1]) != 0: //Запрещенная последовательность "разделитель->разделитель"
				correctexpression = false
				//fmt.Printf(`wrong sequence "multiple separators together": starting from char %d `, index)
			case IsParenthesis(Expression[index]) != 0 && IsSeparator(Expression[index+1]) != 0: //Запрещенная последовательность "скобка->разделитель дроби"
				correctexpression = false
				//fmt.Printf(`wrong sequence "parenthesis->separator": chars %d, %d `, index, index+1)
			case IsParenthesis(Expression[index+1]) != 0 && IsSeparator(Expression[index]) != 0: //Запрещенная последовательность "разделитель дроби->скобка"
				correctexpression = false
				//fmt.Printf(`wrong sequence "separator->parenthesis": chars %d, %d `, index, index+1)
			case IsSeparator(Expression[index]) != 0 && IsOperator(Expression[index+1]) != 0: //Запрещенная последовательность "разделитель дроби->оператор
				correctexpression = false
				//fmt.Printf(`wrong sequence "separator->operation sign": chars %d, %d `, index, index+1)
			case IsSeparator(Expression[index+1]) != 0 && IsOperator(Expression[index]) != 0: //Запрещенная последовательность "оператор->разделитель дроби"
				correctexpression = false
				//fmt.Printf(`wrong sequence "operation sign->separator": chars %d, %d `, index, index+1)
			case IsSeparator(Expression[index]) != 0 && IsNumber(Expression[index+1]) && IsNumber(Expression[index-1]): //Запрещенная последовательность "множественные разделители дроби в числе"
				for nextcharindex := index + 1; nextcharindex < expressionlength; nextcharindex++ {
					if !IsNumber(Expression[nextcharindex]) {
						if IsSeparator(Expression[nextcharindex]) != 0 {
							correctexpression = false
							//fmt.Printf(`wrong sequence "multiple separators within number": starting from char %d `, index)
							break
						} else {
							break
						}
					}
				}
			case IsParenthesis(Expression[index]) == IsLeftParenthesis && IsParenthesis(Expression[index+1]) == IsRightParenthesis: //Запрещенная последовательность "пустые скобки"
				correctexpression = false
				//fmt.Printf(`wrong sequence "empty parentheses": chars %d, %d `, index, index+1)
			case IsParenthesis(Expression[index]) == IsRightParenthesis && countleftparenthesis == 0: // Запрещенная последовательность "подвыражение начинается с правой скобки"
				countrightparenthesis++
				correctexpression = false
				//fmt.Printf(`wrong sequence "beginning form right parenthesis": on char %d `, index)
			case IsParenthesis(Expression[index]) == IsLeftParenthesis && countleftparenthesis == 0: // Считаем левые и правые скобки
				countleftparenthesis++
				for nextcharindex := index + 1; nextcharindex < expressionlength; nextcharindex++ {
					if IsParenthesis(Expression[nextcharindex]) == IsLeftParenthesis {
						countleftparenthesis++
					} else if IsParenthesis(Expression[nextcharindex]) == IsRightParenthesis {
						countrightparenthesis++
					}

				}
			}
		} else if !IsNumber(Expression[index]) && IsParenthesis(Expression[index]) == 0 && IsOperator(Expression[index]) == 0 && IsSeparator(Expression[index]) == 0 { //Недопустимые символы
			correctexpression = false
			//fmt.Printf("incorrect symbol, char %d. Allowed only: %s", index, "1234567890.*/+-()")
		} else if !IsNumber(Expression[index]) && index == expressionlength-1 {
			correctexpression = false
			//fmt.Printf(`wrong sequence "non-numeric last character"`)
		}
	}

	if countleftparenthesis < countrightparenthesis { // Не хватает левых скобок
		correctexpression = false
		//fmt.Printf(`wrong sequence "insufficient number of left parentheses"`)
	} else if countleftparenthesis > countrightparenthesis { // Не хватает правых скобок
		correctexpression = false
		//fmt.Printf(`wrong sequence "insufficient number of right parentheses"`)
	}

	if !correctexpression {
		return false, fmt.Errorf("incorrect expression")
	}
	return true, nil
}

func TokenizeandCalc(Expression string) (float64, error) {
	var result float64
	var operatorsslice []int
	var numsslice []float64
	var priority, tempop, countdown int
	var diverr error

	length := len(Expression)
	for indexoftokenizer := 0; indexoftokenizer < length; indexoftokenizer++ {
		operatorslicelength := len(operatorsslice)
		numsliceslength := len(numsslice)
		if IsNumber(Expression[indexoftokenizer]) {
			num, multiplyindex := ExtractNum(Expression, indexoftokenizer)
			indexoftokenizer = multiplyindex
			numinfloat, _ := strconv.ParseFloat(num, 8)
			numsslice = append(numsslice, numinfloat)
		}
		if !IsNumber(Expression[indexoftokenizer]) && IsSeparator(Expression[indexoftokenizer]) == 0 {
			switch {
			case IsParenthesis(Expression[indexoftokenizer]) == 0 && operatorslicelength == 0:
				operatorsslice = append(operatorsslice, IsOperator(Expression[indexoftokenizer]))

			case IsOperator(Expression[indexoftokenizer]) != 0:
				priority = GetPryority(IsOperator(Expression[indexoftokenizer]))
				if operatorslicelength-1 >= 0 {
					if GetPryority(operatorsslice[operatorslicelength-1]) == priority {
						tempop = operatorsslice[operatorslicelength-1]
						result, diverr = MathOp(numsslice, tempop)
						if diverr != nil {
							return 0, diverr

						}
						numsslice, _ = PopNum(numsslice, len(numsslice)-1)
						if len(numsslice) > 0 {
							numsslice, _ = PopNum(numsslice, len(numsslice)-1)
						}
						if len(operatorsslice) > 0 {
							operatorsslice, _ = PopOp(operatorsslice, len(operatorsslice)-1)
						}
						numsslice = append(numsslice, result)
						operatorsslice = append(operatorsslice, IsOperator(Expression[indexoftokenizer]))
					} else if GetPryority(operatorsslice[operatorslicelength-1]) < priority {
						operatorsslice = append(operatorsslice, IsOperator(Expression[indexoftokenizer]))
					} else if GetPryority(operatorsslice[operatorslicelength-1]) > priority {
						tempop = operatorsslice[operatorslicelength-1]
						result, _ = MathOp(numsslice, tempop)
						if diverr != nil {
							return 0, diverr

						}
						numsslice, _ = PopNum(numsslice, numsliceslength-1)
						numsslice, _ = PopNum(numsslice, numsliceslength-1)
						operatorsslice, _ = PopOp(operatorsslice, operatorslicelength-1)
						operatorsslice = append(operatorsslice, IsOperator(Expression[indexoftokenizer]))
						numsslice = append(numsslice, result)
					}

				}
			case IsParenthesis(Expression[indexoftokenizer]) == IsLeftParenthesis:
				operatorsslice = append(operatorsslice, IsLeftParenthesis)
			case IsParenthesis(Expression[indexoftokenizer]) == IsRightParenthesis:
				for {
					if (operatorsslice[len(operatorsslice)-1]) == IsLeftParenthesis {
						operatorsslice, _ = PopOp(operatorsslice, len(operatorsslice)-1)
						break
					}
					result, _ = MathOp(numsslice, operatorsslice[len(operatorsslice)-1])
					if diverr != nil {
						return 0, diverr

					}
					numsslice, _ = PopNum(numsslice, len(numsslice)-1)
					numsslice, _ = PopNum(numsslice, len(numsslice)-1)
					operatorsslice, _ = PopOp(operatorsslice, len(operatorsslice)-1)
					numsslice = append(numsslice, result)
				}
			case IsParenthesis(Expression[indexoftokenizer]) == 0 && operatorslicelength > 0:
				operatorsslice = append(operatorsslice, IsOperator(Expression[indexoftokenizer]))
			}
		}
		if indexoftokenizer == length {
			break
		}
	}

	countdown = len(operatorsslice) - 1
	for {
		if countdown < 0 {
			break
		} else {
			result, _ = MathOp(numsslice, operatorsslice[countdown])
			if diverr != nil {
				return 0, diverr

			}
			numsslice, _ = PopNum(numsslice, len(numsslice)-1)
			numsslice, _ = PopNum(numsslice, len(numsslice)-1)
			operatorsslice, _ = PopOp(operatorsslice, len(operatorsslice)-1)
			numsslice = append(numsslice, result)
		}
		countdown--

	}

	return result, nil
}

func Calc(Expression string) (float64, error) {

	check, err := IsCorrectExpression(Expression)
	if check && err == nil {
		return TokenizeandCalc(Expression)
	} else {
		return 0, err
	}

}

func main() {

	fmt.Println(Calc("1+1+1-2*(28-1)-1"))
}
