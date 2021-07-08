package mathsubject

import (
	"fmt"
	"math/rand"
	"strconv"
)

func GetSubNew(difficulty, max, no string) []string {
	noi, _ := strconv.Atoi(no)
	myMax, _ := strconv.Atoi(max)
	myDifficulty, _ := strconv.Atoi(difficulty)
	var m = make(map[string]struct{})
	for {
		r := getOne(myDifficulty, myMax)
		_, ok := m[r]
		if !ok {
			m[r] = struct{}{}
			if len(m) > noi-1 {
				break
			}
		}
	}
	s := make([]string, 0, len(m))
	var c = '①'
	var count rune = 0
	for k := range m {
		s = append(s, string(count/3+c)+" "+k)
		count += 1
	}
	return s
}

func getOp(difficulty int) string {
	var op string
	switch difficulty {
	case 1, 4:
		op = "+"
	case 2, 5:
		op = "-"
	default:
		if rand.Intn(2) == 0 {
			op = "+"
		} else {
			op = "-"
		}
	}
	return op
}


func getOne(difficulty, max int) string {
	var r string
	var myDif int
	if difficulty == 7 {
		myDif = rand.Intn(difficulty) + 1
	} else {
		myDif = difficulty
	}
	switch myDif + 1 {
	case 1, 2, 3:
		var op = getOp(difficulty)
		var one = rand.Intn(max)
		var two int
		if "+" == op {
			two = rand.Intn(max - one)
		} else {
			two = rand.Intn(one + 1)
		}
		var fm = "%-2d %s %-2d ="
		r = fmt.Sprintf(fm, one, op, two)
	default:
		var op1 = getOp(difficulty)
		var one = rand.Intn(max)
		var two int
		var three int
		var maxTemp int
		if "+" == op1 {
			two = rand.Intn(max - one)
			maxTemp = one + two
		} else {
			two = rand.Intn(one + 1)
			maxTemp = one - two
		}
		var op2 = getOp(difficulty)
		if "+" == op2 {
			three = rand.Intn(max - maxTemp)
		} else {
			three = rand.Intn(maxTemp + 1)
		}
		var fm = "%-2d%s%-2d%s%-2d="
		r = fmt.Sprintf(fm, one, op1, two, op2, three)
	}
	return r
}
