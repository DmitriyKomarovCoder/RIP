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
	query := ctx.DefaultQuery("query", "")
	if query != "" {
		searchResults := []ds.Company{}

		for _, company := range *companysList {
			if strings.Contains(strings.ToLower(company.CompanyName), strings.ToLower(query)) {
				searchResults = append(searchResults, company)
			}
		}

		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"card":  searchResults,
			"query": query,
		})
	} else {
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"card":  companysList,
			"query": query,
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
