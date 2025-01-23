package handler

import (
	"net/http"
	"strconv"

	"github.com/Vladislave00/CashFlow/app/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createTransaction(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accounting_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	// account_id, err := strconv.Atoi(c.Param("accountid"))
	// if err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, "invalid account id param")
	// 	return
	// }

	var input models.Transaction

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Transaction.CreateTransaction(userId, accounting_id, input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getTransactions(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accounting_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	transactions, err := h.services.Transaction.GetAll(userId, accounting_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *Handler) getTransactionsByAccountId(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accounting_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	account_id, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	transactions, err := h.services.Transaction.GetByAccountId(userId, accounting_id, account_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *Handler) getTransactionById(c *gin.Context) {
	id, err := GetUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	transaction_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	transaction, err := h.services.Transaction.GetById(id, transaction_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *Handler) updateTransaction(c *gin.Context) {

}

func (h *Handler) deleteTransaction(c *gin.Context) {
	id, err := GetUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	transaction_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.Transaction.Delete(id, transaction_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
