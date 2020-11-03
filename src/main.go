package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"math/rand"
	"net/http"
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
		var fm = "%-2d %s %-2d ="
		var r = fmt.Sprintf(fm, one, op, two)
		//var r = strconv.Itoa(one) + " " + op + " " + strconv.Itoa(two) + " = "
		_, ok := m[r]
		if !ok {
			m[r] = struct{}{}
			if len(m) > 50 {
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

func main() {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }, "add": func(i, j int) int { return i + j }})
	r.LoadHTMLGlob("templates/*")
	v1 := r.Group("")
	{
		v1.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", gin.H{
				"title": "练习",
				"items": getSub(),
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
	r.Run(":8081")
}
