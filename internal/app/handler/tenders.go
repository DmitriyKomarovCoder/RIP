package handler

import (
	"RIP/internal/app/ds"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getTenders(ctx *gin.Context) {
	tendersList, err := h.Repository.GetOpenTenders()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tenders"})
		return
	}
	query := ctx.DefaultQuery("query", "")
	if query != "" {
		searchResults := []ds.Tenders{}

		for _, tender := range *tendersList {
			if strings.Contains(strings.ToLower(tender.TenderName), strings.ToLower(query)) {
				searchResults = append(searchResults, tender)
			}
		}

		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"card": searchResults,
		})
	} else {
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"card": tendersList,
		})
	}

}

func (h *Handler) getTenderDetails(ctx *gin.Context) {
	idGet := ctx.Param("id")
	id, _ := strconv.Atoi(idGet)
	tender, err := h.Repository.GetTenderById(id)
	if err != nil {
		ctx.String(http.StatusNotFound, "404 - Not Found")
		return
	}
	ctx.HTML(http.StatusOK, "pages.tmpl", gin.H{
		"card": tender,
	})
}

func (h *Handler) postTenderStatus(ctx *gin.Context) {
	tenderID := ctx.PostForm("tender_id")
	h.Repository.DeleteTender(tenderID)
	ctx.Redirect(http.StatusFound, "/tenders")
}
