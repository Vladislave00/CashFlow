package handler

import (
	"net/http"
	"strconv"

	"github.com/Vladislave00/CashFlow/app/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createAccounting(c *gin.Context) {
	id, err := GetUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input models.Accounting

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accounting_id, err := h.services.Accounting.CreateAccounting(id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": accounting_id,
	})
}

type getAccountingsResponse struct {
	Data []models.Accounting `json:"data"`
}

func (h *Handler) getAccountings(c *gin.Context) {
	id, err := GetUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accountings, err := h.services.Accounting.GetAll(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAccountingsResponse{
		Data: accountings,
	})
}

func (h *Handler) getAccountingById(c *gin.Context) {
	id, err := GetUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accounting_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	accounting, err := h.services.Accounting.GetById(id, accounting_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounting)
}

func (h *Handler) updateAccounting(c *gin.Context) {
	id, err := GetUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accounting_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.UpdateAccountingInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Accounting.Update(id, accounting_id, *input.Name); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteAccounting(c *gin.Context) {
	id, err := GetUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accounting_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.Accounting.Delete(id, accounting_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
