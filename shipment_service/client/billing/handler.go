package billing

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	BillingConnection *BillingConnectionAdapter
}

func NewHandler() *Handler {
	return &Handler{
		BillingConnection: &BillingConnectionAdapter{},
	}
}

func (h *Handler) CreateInvoice(ctx *gin.Context) {

}
