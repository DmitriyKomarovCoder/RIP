package handler

import (
	"RIP/internal/app/ds"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) TenderList(ctx *gin.Context) {
	queryStatus, _ := ctx.GetQuery("status")

	queryStart, _ := ctx.GetQuery("start")

	queryEnd, _ := ctx.GetQuery("end")

	tenders, err := h.Repository.TenderList(queryStatus, queryStart, queryEnd, ctx.GetInt(userCtx), ctx.GetBool(adminCtx))

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"tenders": tenders})
}

func (h *Handler) GetTenderById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	req, com, err := h.Repository.GetTenderWithDataByID(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tender":    req,
		"companies": com,
	})
}

func (h *Handler) UpdateTender(ctx *gin.Context) {
	var updatedTender ds.Tender
	if err := ctx.BindJSON(&updatedTender); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if updatedTender.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("id некоректен"))
		return
	}
	if err := h.Repository.UpdateTender(&updatedTender); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "")
}

func (h *Handler) CreateDraft(c *gin.Context) {
	draftID, err := h.Repository.CreateTenderDraft(creatorID)

	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{"draftID": draftID})
}

func (h *Handler) FormTenderRequest(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.Repository.FormTenderRequestByID(uint(id), creatorID)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	req, com, err := h.Repository.GetTenderWithDataByID(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tender":    req,
		"companies": com,
	})
}

func (h *Handler) RejectTenderRequest(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.Repository.RejectTenderRequestByID(uint(id), moderatorID); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, "отклонена")
}

func (h *Handler) FinishTenderRequest(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.Repository.FinishEncryptDecryptRequestByID(uint(id), moderatorID); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, "завершена")
}
func (h *Handler) DeleteCompanyFromRequest(c *gin.Context) {
	var deleteFromTender ds.TenderCompany
	if err := c.BindJSON(&deleteFromTender); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}
	if deleteFromTender.TenderID <= 0 {
		h.errorHandler(c, http.StatusBadRequest, errors.New("id не найден"))
		return
	}

	if deleteFromTender.CompanyID <= 0 {
		h.errorHandler(c, http.StatusBadRequest, errors.New("id не найден"))
		return
	}

	request, companies, err := h.Repository.DeleteCompanyFromRequest(deleteFromTender)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Компания удалена из заявки", "companies": companies, "monitoring-request": request})
}

func (h *Handler) DeleteTender(c *gin.Context) {
	//ModeratorID тут проверка, что это модератор -> 5 lab
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.Repository.DeleteTenderByID(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "deleted")
}

func (h *Handler) UpdateTenderCompany(c *gin.Context) {
	var TenderCompany ds.TenderCompany
	if err := c.BindJSON(&TenderCompany); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if TenderCompany.TenderID == 0 || TenderCompany.CompanyID == 0 {
		h.errorHandler(c, http.StatusBadRequest, errors.New("не верные id тендера или кампапии"))
		return
	}

	err := h.Repository.UpdateTenderCompany(TenderCompany.TenderID, TenderCompany.CompanyID, TenderCompany.Cash)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, "update")
}
