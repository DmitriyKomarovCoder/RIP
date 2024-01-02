package handler

import (
	"RIP/internal/app/ds"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// TenderList godoc
// @Summary      Get list of tender requests
// @Description  Retrieves a list of tender requests based on the provided parameters
// @Tags         Tenders
// @Accept       json
// @Produce      json
// @Param        status      query  string    false  "Tender request status"
// @Param        start  query  string    false  "Start date in the format '2006-01-02T15:04:05Z'"
// @Param        end    query  string    false  "End date in the format '2006-01-02T15:04:05Z'"
// @Success      200  {object}  []ds.Tender
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router       /api/tenders [get]
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

// GetTenderById godoc
// @Summary      Get tender request by ID
// @Description  Retrieves a tender request with the given ID
// @Tags         Tenders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Tender Request ID"
// @Success      200  {object}  ds.TenderDetails
// @Failure      400  {object}  error
// @Router       /api/tenders/{id} [get]
func (h *Handler) GetTenderById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	req, com, err := h.Repository.GetTenderWithDataByID(uint(id), uint(c.GetInt(userCtx)), c.GetBool(adminCtx))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	tenderD := ds.TenderDetails{Tender: &req, Company: &com}
	c.JSON(http.StatusOK, tenderD)
}

// UpdateTender godoc
// @Summary      Update Tender by admin
// @Description  Update Tender by admin
// @Tags         Tenders
// @Accept       json
// @Produce      json
// @Param        input    body    ds.Tender  true    "updated Assembly"
// @Success      200          {object}  nil
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/tenders [put]
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

//func (h *Handler) CreateDraft(c *gin.Context) {
//	draftID, err := h.Repository.CreateTenderDraft(creatorID)
//
//	if err != nil {
//		h.errorHandler(c, http.StatusInternalServerError, err)
//	}
//
//	c.JSON(http.StatusOK, gin.H{"draftID": draftID})
//}

// FormTenderRequest godoc
// @Summary      Form Company by client
// @Description  Form Company by client
// @Tags         Tenders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Tender form ID"
// @Success      200          {object}  ds.TenderDetails
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/tenders/form/{id} [put]
func (h *Handler) FormTenderRequest(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.Repository.FormTenderRequestByID(uint(id), uint(c.GetInt(userCtx)))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	req, com, err := h.Repository.GetTenderWithDataByID(uint(id), uint(c.GetInt(userCtx)), false)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	tenderDetails := ds.TenderDetails{Tender: &req, Company: &com}
	c.JSON(http.StatusOK, tenderDetails)
}

// UpdateStatusTenderRequest godoc
// @Summary      Update transaction request status by ID
// @Description  Updates the status of a transaction request with the given ID on "завершен"/"отклонен"
// @Tags         Tenders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Request ID"
// @Param        input    body    ds.NewStatus  true    "update status"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /tenders/updateStatus/{id} [put]
func (h *Handler) UpdateStatusTenderRequest(c *gin.Context) {
	var status ds.NewStatus
	if err := c.BindJSON(&status); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if status.Status != "отклонен" && status.Status != "завершен" {
		h.errorHandler(c, http.StatusBadRequest, errors.New("статус можно поменять только на 'отклонен' и 'завершен'"))
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.Repository.FinishRejectHelper(status.Status, uint(id), uint(c.GetInt(userCtx))); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, status.Status+"а")
}

//func (h *Handler) FinishTenderRequest(c *gin.Context) {
//	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
//
//	if err := h.Repository.FinishRejectHelper(uint(id), moderatorID); err != nil {
//		h.errorHandler(c, http.StatusBadRequest, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, "завершена")
//}

// DeleteCompanyFromRequest godoc
// @Summary      Delete company from request
// @Description  Deletes a company from a request based on the user ID and company ID
// @Tags         Tender_Company
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "company ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  error
// @Router       /api/tender-request-company/{id} [delete]
func (h *Handler) DeleteCompanyFromRequest(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	userId := c.GetInt(userCtx)

	request, companies, err := h.Repository.DeleteCompanyFromRequest(uint(userId), uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Компания удалена из заявки", "companies": companies, "monitoring-request": request})
}

// DeleteTender godoc
// @Summary      Delete tender request by user ID
// @Description  Deletes a tender request for the given user ID
// @Tags         Tenders
// @Accept       json
// @Produce      json
// @Param        user_id  path  int  true  "User ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /api/tenders [delete]
func (h *Handler) DeleteTender(c *gin.Context) {
	//userId := c.GetInt(userCtx)
	userId := c.GetInt(userCtx)

	//id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.Repository.DeleteTenderByID(uint(userId))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, "deleted")
}

// UpdateTenderCompany godoc
// @Summary      Update money Tender Company
// @Description  Update money Tender Company by client
// @Tags         Tender_Company
// @Accept       json
// @Produce      json
// @Param        input    	  body    ds.TenderCompany true    "Update money Tender Company"
// @Success      200          {object} map[string]string "update"
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/tender-request-company [put]
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
