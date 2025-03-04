package main

import (
	"fmt"
	"strconv"
	"errors"
)


const isLeftParenthesis = 1
const isRightParenthesis = 2
const isNotParenthesis = 0
const isMultiplication = 10
const isDivision = 20
const isAddition = 30
const isSubtraction = 40
const isNotOperation = 0
const isPoint = 100
const isNotSeparator = 0

var (
	EmptyExpressionErr     = errors.New(`empty expression`)
	IncorrectExpressionErr = errors.New(`incorrect expression`)
	NumToPopMErr           = errors.New(`numtopop > length of slice of nums`)
	NumToPopZeroErr        = errors.New(`numtopop <= 0`)
	NthToPopErr            = errors.New(`no operator to pop`)
	DvsByZeroErr           = errors.New(`division by zero`)
)


func popNum(sliceofnums []float64, numtopop int) ([]float64, []float64, error) {

	var poppednum, newsliceofnums []float64

	if numtopop > len(sliceofnums) {
		return poppednum, sliceofnums, NumToPopMErr // NumToPopMErr
	}
	if numtopop <= 0 {
		return poppednum, sliceofnums, NumToPopMErr // NumToPopZeroErr
	}

	poppednum = append(sliceofnums[len(sliceofnums)-numtopop:])
	newsliceofnums = append(sliceofnums[:len(sliceofnums)-numtopop], sliceofnums[len(sliceofnums):]...)

	return poppednum, newsliceofnums, nil
}

func popOp(opslice []int) (int, []int, error) {

	var newopslice []int

	if len(opslice) == 0 {
		return 0, opslice, NthToPopErr // NthToPopErr
	}

	poppedop := opslice[len(opslice)-1]
	newopslice = append(opslice[:len(opslice)-1], opslice[len(opslice):]...)

	return poppedop, newopslice, nil
}


func GetLightExpressions(sliceofnums []float64, opslice []int, operator int, prioritynum int, addop bool) ([]float64, map[int]string, error) {
	var poppedop int
	var poppednums []float64
	var popnumerr, popoperr error
	lightexprs := make(map[int]string)
	poppedop, opslice, popoperr = popOp(opslice)
	poppednums, sliceofnums, popnumerr = popNum(sliceofnums, 2)

	if poppedop == 0 && popoperr != nil {
		return sliceofnums, map[int]string{}, popoperr
	}

	if poppednums == nil && popnumerr != nil {
		return sliceofnums, map[int]string{}, popnumerr
	}

	switch {
	case poppedop == isAddition:
		fnum := strconv.FormatFloat(poppednums[0], 'g', 8, 64)
		snum := strconv.FormatFloat(poppednums[1], 'g', 8, 64)
		str := fnum + "+" + snum
		lightexprs[prioritynum] = str
		
	case poppedop == isSubtraction:
		fnum := strconv.FormatFloat(poppednums[0], 'g', 8, 64)
		snum := strconv.FormatFloat(poppednums[1], 'g', 8, 64)
		str := fnum + "-" + snum
		lightexprs[prioritynum] = str

	case poppedop == isMultiplication:
		fnum := strconv.FormatFloat(poppednums[0], 'g', 8, 64)
		snum := strconv.FormatFloat(poppednums[1], 'g', 8, 64)
		str := fnum + "*" + snum
		lightexprs[prioritynum] = str

	case poppedop == isDivision:
		if poppednums[1] == 0 {
			return sliceofnums, lightexprs, DvsByZeroErr //DvsByZeroErr
		}
		fnum := strconv.FormatFloat(poppednums[0], 'g', 8, 64)
		snum := strconv.FormatFloat(poppednums[1], 'g', 8, 64)
		str := fnum + "/" + snum
		lightexprs[prioritynum] = str
	}
	return sliceofnums, lightexprs, nil
}

/*
const s = "-+/*"
const isOperator = 1
const otherSymb = 0

func getSymb(char byte) int {
	for i := range s {
		if string(char) == string(s[i]) {
			return isOperator
		}
	}
	return otherSymb
}

func GetLightExpressions(Expression string) []string {
	var (
		exprsl []string
	)
	for end := 0; end < 1; {
		for n := 0; n < len(Expression); n++ {
			if getSymb(Expression[n]) == otherSymb {
				if n == len(Expression)-1 {
					exprsl = append(exprsl, Expression[:n+1])
					end++
					return exprsl
				}
				continue
			} else {
				if n+1 == len(Expression)-1 {
					exprsl = append(exprsl, Expression[1:])
					return exprsl
				}
				for k := n + 1; getSymb(Expression[k]) != isOperator && k < len(Expression) && n < len(Expression); k++ {
					n++
				}
				//fmt.Println(Expression)
				exprsl = append(exprsl, Expression[:n+1])

				Expression = Expression[n+1:]
				//fmt.Println(Expression)
				n = 0

			}
		}
	}
	return exprsl
}
*/

func main() {
	str := "1+2+1+1+(1+1)+1+1"
	/*for i := range sl {
		fmt.Println(string(sl[i]))
	}
	fmt.Println(sl)
	*/

	GetLightExpressions()

	fmt.Println()
}
