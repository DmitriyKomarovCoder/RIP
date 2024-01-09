package handler

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/role"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func ParseDateString(dateString string) (time.Time, error) {
	format := "2006-01-02"
	parsedTime, err := time.Parse(format, dateString)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

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
	userID, existsUser := ctx.Get("user_id")
	userRole, existsRole := ctx.Get("user_role")
	if !existsUser || !existsRole {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}

	switch userRole {
	case role.Buyer:
		h.tenderByUserId(ctx, fmt.Sprintf("%d", userID))
		return
	default:
		break
	}

	queryStatus := ctx.Query("status_id")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" {
		startDateStr = "0001-01-01"
	}
	if endDateStr == "" {
		endDateStr = time.Now().Format("2006-01-02")
	}

	startDate, errStart := ParseDateString(startDateStr)
	endDate, errEnd := ParseDateString(endDateStr)
	h.Logger.Info(startDate, endDate)
	if errEnd != nil || errStart != nil {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("incorrect `start_date` or `end_date`"))
		return
	}

	tenders, err := h.Repository.TendersList(queryStatus, startDate, endDate)

	if err != nil {
		h.errorHandler(ctx, http.StatusNoContent, err)
		return
	}
	h.successHandler(ctx, "tenders", tenders)
}

func (h *Handler) tenderByUserId(ctx *gin.Context, userID string) {
	tenders, errDB := h.Repository.TenderByUserID(userID)
	if errDB != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, errDB)
		return
	}

	h.successHandler(ctx, "tenders", tenders)
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

	//req, com, err := h.Repository.GetTenderWithDataByID(uint(id), uint(c.GetInt(userCtx)), c.GetBool(adminCtx))
	tender, err := h.Repository.TenderByID(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	//tenderD := ds.TenderDetails{Tender: &req, Company: &com}
	h.successHandler(c, "tender", tender)
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
	userID, existsUser := ctx.Get("user_id")
	userRole, existsRole := ctx.Get("user_role")
	if !existsUser || !existsRole {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}

	var updatedTender ds.UpdateTender
	if err := ctx.BindJSON(&updatedTender); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if updatedTender.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("id некоректен"))
		return
	}

	var updatedT ds.Tender
	updatedT.ID = updatedTender.ID
	updatedT.Name = updatedTender.Name

	tender, err := h.Repository.TenderByID(updatedT.ID)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("hike with `id` = %d not found", tender.ID))
		return
	}

	if tender.UserID != userID && userRole == role.Buyer {
		h.errorHandler(ctx, http.StatusForbidden, errors.New("you cannot change the hike if it's not yours"))
		return
	}

	if err := h.Repository.UpdateTender(&updatedT); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "updated_tender", gin.H{
		"id":              updatedTender.ID,
		"tender_name":     updatedTender.Name,
		"creation_date":   tender.CreationDate,
		"completion_date": tender.CompletionDate,
		"formation_date":  tender.FormationDate,
		"user_id":         tender.UserID,
		"status":          tender.Status,
	})
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
// @Router       /api/tenders/form [put]
func (h *Handler) FormTenderRequest(c *gin.Context) {
	//id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID, existsUser := c.Get("user_id")
	if !existsUser {
		h.errorHandler(c, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}

	err, _ := h.Repository.FormTenderRequestByID(userID.(uint))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	//_, _, err = h.Repository.GetTenderWithDataByID(idTender, userID.(uint), false)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	//tenderDetails := ds.TenderDetails{Tender: &req, Company: &com}
	c.Status(http.StatusOK)
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
// @Router       /tenders/updateStatus [put]
func (h *Handler) UpdateStatusTenderRequest(c *gin.Context) {
	var status ds.NewStatus
	if err := c.BindJSON(&status); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	userIDStr, existsUser := c.Get("user_id")
	if !existsUser {
		h.errorHandler(c, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}
	userID := userIDStr.(uint)

	if status.Status != "отклонен" && status.Status != "завершен" {
		h.errorHandler(c, http.StatusBadRequest, errors.New("статус можно поменять только на 'отклонен' и 'завершен'"))
	}

	//id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.Repository.FinishRejectHelper(status.Status, status.TenderID, userID); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
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
// @Router       /api/tender-request-company [delete]
func (h *Handler) DeleteCompanyFromRequest(c *gin.Context) {
	var body struct {
		ID int `json:"id"`
	}

	if err := c.BindJSON(&body); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if body.ID == 0 {
		h.errorHandler(c, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}

	err := h.Repository.DeleteCompanyFromRequest(body.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	h.successHandler(c, "deleted_company_tender", body.ID)
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
	userID, existsUser := c.Get("user_id")
	userRole, existsRole := c.Get("user_role")
	if !existsUser || !existsRole {
		h.errorHandler(c, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}

	//userId := c.GetInt(userCtx)
	var request struct {
		ID uint `json:"id"`
	}

	if err := c.BindJSON(&request); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if request.ID == 0 {
		h.errorHandler(c, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}

	//userId := c.GetInt(userCtx)

	tender, err := h.Repository.TenderByID(request.ID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, fmt.Errorf("tender with `id` = %d not found", tender.ID))
		return
	}

	if tender.UserID != userID && userRole == role.Buyer {
		h.errorHandler(c, http.StatusForbidden, errors.New("you are not the creator. you can't delete a tender"))
		return
	}

	err = h.Repository.DeleteTenderByID(request.ID)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	h.successHandler(c, "tender_id", request.ID)
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
	//var TenderCompany ds.TenderCompany
	var TenderCompanyU ds.TenderCompanyUpdate
	if err := c.BindJSON(&TenderCompanyU); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	//if TenderCompanyU.TenderID == 0 || TenderCompanyU.CompanyID == 0 {
	//	h.errorHandler(c, http.StatusBadRequest, errors.New("не верные id тендера или кампапии"))
	//	return
	//}

	err := h.Repository.UpdateTenderCompany(TenderCompanyU.ID, TenderCompanyU.Cash)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, "update")
}
