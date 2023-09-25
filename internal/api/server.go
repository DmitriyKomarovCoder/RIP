package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Company struct {
	Id          int
	Url         string
	Img         string
	Title       string
	Description string
	IIN         string
	Income      int
}

func StartServer() {

	companySlice := []Company{
		{1, "company/1", "image/img1.png", "Компания 1", `Повседневная практика показывает, что реализация намеченных плановых заданий обеспечивает широкому кругу (специалистов) участие в формировании направлений прогрессивного развития. Не следует, однако забывать, что дальнейшее развитие различных форм деятельности требуют определения и уточнения дальнейших направлений развития.

		Идейные соображения высшего порядка, а также постоянный количественный рост и сфера нашей активности требуют определения и уточнения системы обучения кадров, соответствует насущным потребностям. С другой стороны укрепление и развитие структуры требуют определения и уточнения дальнейших направлений развития.`, "3449013711 / 324901001", 750000000},
		{2, "company/2", "image/img2.png", "Компания 2", `Таким образом постоянный количественный рост и сфера нашей активности в значительной степени обуславливает создание существенных финансовых и административных условий. Не следует, однако забывать, что сложившаяся структура организации обеспечивает широкому кругу (специалистов) участие в формировании дальнейших направлений развития. Повседневная практика показывает, что реализация намеченных плановых заданий влечет за собой процесс внедрения и модернизации дальнейших направлений развития. `, "2449013711 / 2434901001", 23423422},
		{3, "company/3", "image/img3.jpg", "Компания 3", `Не следует, однако забывать, что начало повседневной работы по формированию позиции требуют определения и уточнения соответствующий условий активизации. Значимость этих проблем настолько очевидна, что дальнейшее развитие различных форм деятельности способствует подготовки и реализации форм развития. Идейные соображения высшего порядка, а также сложившаяся структура организации позволяет выполнять важные задания по разработке направлений прогрессивного развития. С другой стороны новая модель организационной деятельности способствует подготовки и реализации дальнейших направлений развития. `, "1449013711 / 2134901001", 3242343244},
	}

	log.Println("Server start up")

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/image", "./resources")

	r.GET("/companys", func(c *gin.Context) {
		query := c.DefaultQuery("query", "")
		if query != "" {
			searchResults := []Company{}

			for _, company := range companySlice {
				if strings.Contains(strings.ToLower(company.Title), strings.ToLower(query)) {
					fmt.Println(company.Title, query)
					searchResults = append(searchResults, company)
				}
			}
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"card":    searchResults,
				"CSSFile": "image/card.css",
				"query":   query,
			})

		} else {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"card":    companySlice,
				"CSSFile": "image/card.css",
				"query":   query,
			})
		}
	})

	r.GET("/companys/:id", func(c *gin.Context) {
		idGet := c.Param("id")
		id, _ := strconv.Atoi(idGet)
		if id > len(companySlice) {
			c.String(http.StatusNotFound, "404 - Not Found")
			return
		}

		c.HTML(http.StatusOK, "pages.tmpl", gin.H{
			"card":    companySlice[id-1],
			"CSSFile": "image/card.css",
		})
	})

	r.Run()

	log.Println("Server down")
}
