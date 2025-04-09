package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// todo: Steps:
// * tangkap parameter di handler
// * handler ke service
// * service yang menentukan repository mana yang di-call
// * repository: FindAll, FindByUserID
// * db

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

// note: GET: /api/v1/campaigns
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	// todo: Convert string user_id to int
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaigns)
	c.JSON(http.StatusOK, response)
}
