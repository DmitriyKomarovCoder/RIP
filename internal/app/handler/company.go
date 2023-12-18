package handler

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func (h *Handler) CompaniesList(ctx *gin.Context) {
	queryText, _ := ctx.GetQuery("company_name")
	companies, err := h.Repository.CompaniesList(queryText)
	if err != nil {
		h.errorHandler(ctx, http.StatusNoContent, err)
		return
	}
	draftID, err := h.Repository.GetTenderDraftID(ctx.GetInt(userCtx)) // creatorID(UserID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	companiesList := ds.CompanyList{
		DraftID:   draftID,
		Companies: companies,
	}

	h.successHandler(ctx, "companies", companiesList)
}

func (h *Handler) GetCompanyById(ctx *gin.Context) {
	//queryText, _ := ctx.GetQuery("company_name")

	id, err := strconv.ParseUint(ctx.Param("id")[:], 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
	}

	company, err := h.Repository.GetCompanyById(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "company", company)
}

func (h *Handler) DeleteCompany(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id")[:], 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if id == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}

	url := h.Repository.DeleteCompanyImage(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.DeleteImage(utils.ExtractObjectNameFromUrl(url))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteCompany(uint(id))

	if gorm.IsRecordNotFoundError(err) {
		h.errorHandler(ctx, http.StatusBadRequest, err)
	} else if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, "company deleted successfully")
}

func (h *Handler) AddCompany(ctx *gin.Context) {
	var newCompany ds.Company

	if newCompany.ID != 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}

	newCompany.CompanyName = ctx.Request.FormValue("name")
	if newCompany.CompanyName == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("имя компании не может быть пустой"))
		return
	}

	newCompany.IIN = ctx.Request.FormValue("IIN")
	if newCompany.IIN == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("имя ИИН не может быть пустой"))
		return
	}

	newCompany.Description = ctx.Request.FormValue("description")
	if newCompany.Description == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("описание не может быть пустой"))
		return
	}

	file, header, err := ctx.Request.FormFile("image")
	if err != http.ErrMissingFile && err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ошибка при загрузке изображения"})
		return
	}

	if newCompany.ImageURL, err = h.SaveImage(ctx.Request.Context(), file, header); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ошибка при сохранении изображения"})
		return
	}

	create_id, err := h.Repository.AddCompany(&newCompany)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successAddHandler(ctx, "company_id", create_id)
}

func (h *Handler) UpdateCompany(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id")[:], 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	file, header, err := ctx.Request.FormFile("image")
	// if err != nil {
	// 	h.errorHandler(ctx, http.StatusBadRequest, err)
	// 	return
	// }

	// var updatedCompany ds.Company
	// if err := ctx.BindJSON(&updatedCompany); err != nil {
	// 	h.errorHandler(ctx, http.StatusBadRequest, err)
	// 	return
	// }

	var updatedCompany ds.Company

	updatedCompany.ID = uint(id)

	if updatedCompany.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("param `id` not found"))
	}

	updatedCompany.CompanyName = ctx.Request.FormValue("name")
	updatedCompany.IIN = ctx.Request.FormValue("IIN")
	updatedCompany.Description = ctx.Request.FormValue("description")

	if header != nil && header.Size != 0 {
		if updatedCompany.ImageURL, err = h.SaveImage(ctx.Request.Context(), file, header); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		url := h.Repository.DeleteCompanyImage(updatedCompany.ID)

		if err = h.DeleteImage(utils.ExtractObjectNameFromUrl(url)); err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
	}

	if _, err := h.Repository.UpdateCompany(&updatedCompany); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successHandler(ctx, "updated_company", gin.H{
		"id":           updatedCompany.ID,
		"company_name": updatedCompany.CompanyName,
		"description":  updatedCompany.Description,
		"image_url":    updatedCompany.ImageURL,
		"status":       updatedCompany.Status,
		"iin":          updatedCompany.IIN,
	})
}

func (h *Handler) AddCompanyToRequest(ctx *gin.Context) {
	var request ds.AddToCompanyID

	request.UserID = ctx.GetInt(userCtx)
	idStr := ctx.Param("id")
	request.CompanyID = uint(ctx.GetInt(idStr))

	if request.CompanyID == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "услуга не может быть пустой"})
		return
	}

	draftID, err := h.Repository.AddCompanyToDraft(request.CompanyID, uint(request.UserID))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"draftID": draftID,
	})
}

// func (h *Handler) postCompanyStatus(ctx *gin.Context) {
// 	companyID := ctx.PostForm("company_id")
// 	h.Repository.DeleteCompany(companyID)
// 	ctx.Redirect(http.StatusFound, "/companys")
// }
