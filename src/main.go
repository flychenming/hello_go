package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getSubNew(difficulty, max, no string) []string {
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

func uploadFile(c *gin.Context) {
	// FormFile方法会读取参数“upload”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	// header调用Filename方法，就可以得到文件名
	filename := header.Filename
	fmt.Println(file, err, filename)

	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	out, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	c.String(http.StatusCreated, "upload successful \n")
}

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

func download(c *gin.Context){
	file := c.Query("file")
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))//fmt.Sprintf("attachment; filename=%s", filename) Downloaded file renamed
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(file)
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
				s = append(s, subject{Logic: getRandomText(logic), Item: getSubNew(difficulty, max, no)})
			}
			c.HTML(200, "my.html", gin.H{
				"title": "练习",
				"items": s,
			})
		})
		v1.GET("/upload", func(c *gin.Context) {
			c.HTML(200, "upload.html", gin.H{})
		})
		v1.POST("/upload", uploadFile)

		v1.GET("/download", download)

		v1.GET("/list", func(c *gin.Context) {
			files, err := ioutil.ReadDir(".")
			if err != nil {
				log.Fatal(err)
			}
			l := make([]string, 0)
			for _, f := range files {
				if !f.IsDir() && f.Name()!="Dockerfile" && f.Name()!="hello_go.iml"{
					fmt.Println(f.Name())
					l=append(l, f.Name())
				}
			}
			c.HTML(200, "list.html", gin.H{
				"items":l,
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
