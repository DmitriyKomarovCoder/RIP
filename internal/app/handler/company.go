package handler

import (
	"RIP/internal/app/ds"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getCompanys(ctx *gin.Context) {
	companysList, err := h.Repository.GetOpenCompanys()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tenders"})
		return
	}
	company_name := ctx.DefaultQuery("company_name", "")
	if company_name != "" {
		searchResults := []ds.Company{}

		for _, company := range *companysList {
			if strings.Contains(strings.ToLower(company.CompanyName), strings.ToLower(company_name)) {
				searchResults = append(searchResults, company)
			}
		}

		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"card":         searchResults,
			"company_name": company_name,
		})
	} else {
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"card":         companysList,
			"company_name": company_name,
		})
	}

}

func (h *Handler) getCompanyDetails(ctx *gin.Context) {
	idGet := ctx.Param("id")
	id, _ := strconv.Atoi(idGet)
	tender, err := h.Repository.GetCompanyById(id)
	if err != nil {
		ctx.String(http.StatusNotFound, "404 - Not Found")
		return
	}
	ctx.HTML(http.StatusOK, "pages.tmpl", gin.H{
		"card": tender,
	})
}

func (h *Handler) postCompanyStatus(ctx *gin.Context) {
	companyID := ctx.PostForm("company_id")
	h.Repository.DeleteCompany(companyID)
	ctx.Redirect(http.StatusFound, "/companys")
}
