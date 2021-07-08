package main

import (
	"./mathsubject"
	"./upload"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func getRandomText(max int) []string {
	data, err := ioutil.ReadFile("src/1.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return []string{}
	}
	var content = strings.Split(string(data), "\n")
	s := make([]string, 0, max)
	var c = '①'
	var count rune = 0
	for i := 0; i < max; i++ {
		s = append(s, string(c+count)+content[rand.Intn(len(content))]+"?")
		count += 1
	}
	return s
}

type subject struct {
	Item  []string
	Logic []string
}

func main() {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }, "add": func(i, j int) int { return i + j }})
	r.LoadHTMLGlob("src/templates/*")
	r.Static("/static/css", "src/static/css")
	//r.LoadHTMLGlob(filepath.Join(os.Getenv("GOPATH"), "src/templates/*"))
	v1 := r.Group("")
	{
		v1.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", gin.H{
				"title":      "选择难易度",
				"difficulty": map[int]string{1: "单一加", 2: "单一减", 3: "单一加减混合", 4: "2连加", 5: "2连减", 6: "2连混合", 7: "混合"},
				"max":        []int{10, 20, 30, 50, 100},
				"no":         []int{30, 50},
				"logic":      []int{1, 2, 3, 5, 8},
				"pages":      []int{1, 3, 5, 8, 10},
			})
		})
		v1.GET("/my", func(c *gin.Context) {
			difficulty := c.Query("difficulty")
			max := c.Query("max")
			no := c.Query("no")
			page, _ := strconv.Atoi(c.Query("page"))
			logic, _ := strconv.Atoi(c.Query("logic"))
			s := make([]subject, 0, page)
			for i := 0; i < page; i++ {
				s = append(s, subject{Logic: getRandomText(logic), Item: mathsubject.GetSubNew(difficulty, max, no)})
			}
			c.HTML(200, "my.html", gin.H{
				"title": "练习",
				"items": s,
			})
		})
		v1.GET("/upload", func(c *gin.Context) {
			c.HTML(200, "upload.html", gin.H{})
		})
		v1.POST("/upload", upload.UploadFile)

		v1.GET("/download", upload.Download)

		v1.GET("/list", func(c *gin.Context) {
			files, err := ioutil.ReadDir(".")
			if err != nil {
				log.Fatal(err)
			}
			l := make([]string, 0)
			for _, f := range files {
				if !f.IsDir() && f.Name() != "Dockerfile" && f.Name() != "hello_go.iml" {
					fmt.Println(f.Name())
					l = append(l, f.Name())
				}
			}
			c.HTML(200, "list.html", gin.H{
				"items": l,
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
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}
}
