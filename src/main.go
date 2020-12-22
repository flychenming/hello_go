package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
)

func getSub() []string {
	var m = make(map[string]struct{})
	for ; ; {
		var op string
		if rand.Intn(2) == 0 {
			op = "+"
		} else {
			op = "-"
		}
		var one = rand.Intn(20)
		var two int
		if "+" == op {
			two = rand.Intn(20 - one)
		} else {
			two = rand.Intn(one + 1)
		}
		var fm = "%-2d %s %-2d"
		var r = fmt.Sprintf(fm, one, op, two)
		_, ok := m[r]
		if !ok {
			m[r] = struct{}{}
			if len(m) > 29 {
				break
			}
		}
	}
	s := make([]string, 0, len(m))
	for k, _ := range m {
		s = append(s, k)
	}
	return s
}

func getSubNew(difficulty, max, no string) []string {
	noi, _ := strconv.Atoi(no)
	myMax, _ := strconv.Atoi(max)
	myDifficulty, _ := strconv.Atoi(difficulty)
	var m = make(map[string]struct{})
	for ; ; {
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
	for k, _ := range m {
		s = append(s, k)
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
	switch difficulty {
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

func main() {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }, "add": func(i, j int) int { return i + j }})
	r.LoadHTMLGlob("src/templates/*")
	//r.LoadHTMLGlob(filepath.Join(os.Getenv("GOPATH"), "src/templates/*"))
	v1 := r.Group("")
	{
		v1.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", gin.H{
				"title":      "选择难易度",
				"difficulty": map[int]string{1: "单一加", 2: "单一减", 3: "单一加减混合", 4: "2连加", 5: "2连减", 6: "2连混合"},
				"max":        []int{10, 20, 30, 50, 100},
				"no":         []int{30, 50},
				"items":      getSub(),
			})
		})
		v1.GET("/easy", func(c *gin.Context) {
			c.HTML(200, "easy.html", gin.H{
				"title": "练习",
				"items": getSub(),
			})
		})
		v1.GET("/my", func(c *gin.Context) {
			difficulty := c.Query("difficulty")
			max := c.Query("max")
			no := c.Query("no")
			c.HTML(200, "my.html", gin.H{
				"title": "练习",
				"items": getSubNew(difficulty, max, no),
			})
		})
	}
	//定义默认路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": 404,
			"error":  "404, page not exists!",
		})
	})
	r.Run(":8080")
}
