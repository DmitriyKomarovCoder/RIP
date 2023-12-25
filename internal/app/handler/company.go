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

// CompaniesList godoc
// @Summary      Companies List
// @Description  Companies List
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        name query   string  false  "Query string to filter companies by name"
// @Success      200          {object}  ds.CompanyList
// @Failure      500          {object}  error
// @Router       /api/companies [get]
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

// GetCompanyById godoc
// @Summary      Company By ID
// @Description  Company By ID
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id   path    int     true        "Companies ID"
// @Success      200          {object}  ds.Company
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/companies/{id} [get]
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

// DeleteCompany godoc
// @Summary      Delete company by ID
// @Description  Deletes a company with the given ID
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Company ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /api/companies/{id} [delete]
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

// AddCompany godoc
// @Summary      Add new company
// @Description  Add a new company with image, name, IIN
// @Tags         Companies
// @Accept       multipart/form-data
// @Produce      json
// @Param        image formData file true "Company image"
// @Param        name formData string true "Company name"
// @Param        description formData string false "Company description"
// @Param        IIN formData integer true "Company IIN"
// @Success      201  {string}  map[string]any
// @Failure      400  {object}  map[string]any
// @Router       /api/companies [post]
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

// UpdateCompany godoc
// @Summary      Update company by ID
// @Description  Updates a company with the given ID
// @Tags         Companies
// @Accept       multipart/form-data
// @Produce      json
// @Param        id          path        int     true        "ID"
// @Param        name        formData    string  false       "name"
// @Param        description formData    string  false       "description"
// @Param        IIN         formData    string  false       "IIN"
// @Param        image       formData    file    false       "image"
// @Success      200         {object}    map[string]any
// @Failure      400         {object}    error
// @Router       /api/companies/{id} [put]
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

// AddCompanyToRequest godoc
// @Summary      Add company to request
// @Description  Adds a company to a tender request
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        threatId  path  int  true  "Threat ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /companies/request/{id} [post]
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
