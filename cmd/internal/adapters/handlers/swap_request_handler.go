package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"swapp-go/cmd/internal/adapters/handlers/responses"
	"swapp-go/cmd/internal/application/services"
	"swapp-go/cmd/internal/domain"
)

type SwapRequestHandler struct {
	swapRequestService services.SwapRequestServiceInterface
}

func NewSwapRequestHandler(swapRequestService services.SwapRequestServiceInterface) *SwapRequestHandler {
	return &SwapRequestHandler{swapRequestService}
}

type SwapRequestRequest struct {
	OfferedItemID        uuid.UUID `json:"offered_item_id" binding:"required"`
	RequestedItemID      uuid.UUID `json:"requested_item_id" binding:"required"`
	RequestedItemOwnerID uuid.UUID `json:"requested_item_owner_id" binding:"required"`
}

type SwapRequestResponse struct {
	ID                   string `json:"id"`
	Status               string `json:"status"`
	ReferenceNumber      string `json:"reference_number"`
	OfferedItemID        string `json:"offered_item_id"`
	RequestedItemID      string `json:"requested_item_id"`
	OfferedItemOwnerID   string `json:"offered_item_owner_id"`
	RequestedItemOwnerID string `json:"requested_item_owner_id"`
}

type SwapRequestSuccessResponse struct {
	Message     string               `json:"message"`
	SwapRequest *SwapRequestResponse `json:"swapRequest"`
}

type SwapRequestListResponse struct {
	Message      string                `json:"message"`
	SwapRequests []SwapRequestResponse `json:"swapRequests"`
}

func (handler *SwapRequestHandler) Create(context *gin.Context) {
	var requestInput SwapRequestRequest

	if err := context.ShouldBindJSON(&requestInput); err != nil {
		responses.BadRequest(context, "Invalid request body", err)
		return
	}

	userID, err := getUserIDFromContext(context)
	if err != nil {
		responses.Unauthorized(context, "Unauthorized", err)
		return
	}

	referenceNumber, err := generateReferenceNumber()
	if err != nil {
		fmt.Println(err)
	}

	swapRequest := &domain.SwapRequest{
		Status:               domain.StatusPending,
		ReferenceNumber:      referenceNumber,
		OfferedItemID:        requestInput.OfferedItemID,
		RequestedItemID:      requestInput.RequestedItemID,
		OfferedItemOwnerID:   userID,
		RequestedItemOwnerID: requestInput.RequestedItemOwnerID,
	}

	if err = handler.swapRequestService.Create(swapRequest); err != nil {
		responses.InternalServerError(context, "Failed to create swapp request", err)
		return
	}

	respondWithSwapRequest(context, http.StatusCreated, "Swap request created successfully!", swapRequest)
}

func (handler *SwapRequestHandler) FindByID(context *gin.Context) {
	userID, err := getUserIDFromContext(context)
	if err != nil {
		responses.Unauthorized(context, "Unauthorized", nil)
		return
	}

	requestID, err := uuid.Parse(context.Param("id"))
	if err != nil {
		responses.BadRequest(context, "Invalid ID format", err)
		return
	}

	swapRequest, err := handler.swapRequestService.FindByID(requestID)
	if err != nil {
		responses.NotFound(context, "Swap request not found", err)
		return
	}

	if swapRequest.OfferedItemOwnerID != userID && swapRequest.RequestedItemOwnerID != userID {
		responses.Unauthorized(context, "You are not part of this swap request", nil)
		return
	}

	respondWithSwapRequest(context, http.StatusOK, "Swap request fetched successfully", swapRequest)
}

func (handler *SwapRequestHandler) FindByReferenceNumber(context *gin.Context) {
	userID, err := getUserIDFromContext(context)
	if err != nil {
		responses.Unauthorized(context, "Unauthorized", nil)
		return
	}

	reference := context.Param("reference")
	if reference == "" {
		responses.BadRequest(context, "Missing reference number", nil)
		return
	}

	swapRequest, err := handler.swapRequestService.FindByReferenceNumber(reference)
	if err != nil {
		responses.NotFound(context, "Swap request not found", err)
		return
	}

	if swapRequest.OfferedItemOwnerID != userID && swapRequest.RequestedItemOwnerID != userID {
		responses.Unauthorized(context, "You are not authorized to view this swap request", nil)
		return
	}

	respondWithSwapRequest(context, http.StatusOK, "Swap request retrieved!", swapRequest)
}

func (handler *SwapRequestHandler) Delete(context *gin.Context) {
	userID, err := getUserIDFromContext(context)
	if err != nil {
		responses.Unauthorized(context, "Unauthorized", nil)
		return
	}

	requestID, err := uuid.Parse(context.Param("id"))
	if err != nil {
		responses.BadRequest(context, "Invalid ID format", err)
		return
	}

	swapRequest, err := handler.swapRequestService.FindByID(requestID)
	if err != nil {
		responses.NotFound(context, "Swap request not found", err)
		return
	}

	if swapRequest.OfferedItemOwnerID != userID {
		responses.Unauthorized(context, "Only the sender can delete this request", nil)
		return
	}

	if err = handler.swapRequestService.Delete(requestID); err != nil {
		responses.InternalServerError(context, "Failed to delete request", err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Swap Request deleted successfully!"})
}

func (handler *SwapRequestHandler) UpdateStatus(context *gin.Context) {
	userID, err := getUserIDFromContext(context)
	if err != nil {
		responses.Unauthorized(context, "Unauthorized", nil)
		return
	}

	requestID, err := uuid.Parse(context.Param("id"))
	if err != nil {
		responses.BadRequest(context, "Invalid ID format", err)
		return
	}

	var body struct {
		Status domain.SwapRequestStatus `json:"status"`
	}
	if err = context.ShouldBindJSON(&body); err != nil {
		responses.BadRequest(context, "Invalid status", err)
		return
	}

	swapRequest, err := handler.swapRequestService.FindByID(requestID)
	if err != nil {
		responses.NotFound(context, "Swap request not found", err)
		return
	}

	if swapRequest.OfferedItemOwnerID != userID && body.Status == "cancelled" {
		responses.Unauthorized(context, "Only the sender can cancel this request", nil)
		return
	}

	if swapRequest.RequestedItemOwnerID != userID && body.Status != "cancelled" {
		responses.Unauthorized(context, "Only the recipient can accept or reject a request", nil)
		return
	}

	if err = handler.swapRequestService.UpdateStatus(requestID, body.Status); err != nil {
		responses.InternalServerError(context, "Failed to update status", err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Status updated successfully!"})
}

func (handler *SwapRequestHandler) ListByUser(context *gin.Context) {
	userID, err := getUserIDFromContext(context)
	if err != nil {
		responses.Unauthorized(context, "Unauthorized", nil)
		return
	}

	swapRequests, err := handler.swapRequestService.ListByUser(userID)
	if err != nil {
		responses.InternalServerError(context, "Failed to fetch swap requests", err)
		return
	}

	respondWithSwapRequestList(context, http.StatusOK, "Swap requests fetched successfully", swapRequests)
}

func (handler *SwapRequestHandler) ListByStatus(context *gin.Context) {
	userID, err := getUserIDFromContext(context)
	if err != nil {
		responses.Unauthorized(context, "Unauthorized", nil)
		return
	}

	statusParam := context.Param("status")
	if statusParam == "" {
		responses.BadRequest(context, "Missing status parameter", nil)
		return
	}

	status := domain.SwapRequestStatus(statusParam)

	swapRequests, err := handler.swapRequestService.ListByStatus(status)
	if err != nil {
		responses.InternalServerError(context, "Failed to fetch swap requests", err)
		return
	}

	filtered := make([]domain.SwapRequest, 0)
	for _, request := range swapRequests {
		if request.OfferedItemOwnerID == userID || request.RequestedItemOwnerID == userID {
			filtered = append(filtered, request)
		}
	}

	respondWithSwapRequestList(context, http.StatusOK, "Filtered swap requests by status", filtered)
}

func getUserIDFromContext(context *gin.Context) (uuid.UUID, error) {
	rawUserID, exists := context.Get("userID")

	if !exists {
		return uuid.Nil, errors.New("user not in context")
	}

	return uuid.Parse(rawUserID.(string))
}

func respondWithSwapRequest(context *gin.Context, status int, message string, swapRequest *domain.SwapRequest) {
	response := SwapRequestSuccessResponse{
		Message: message,
		SwapRequest: &SwapRequestResponse{
			ID:                   swapRequest.ID.String(),
			Status:               string(swapRequest.Status),
			ReferenceNumber:      swapRequest.ReferenceNumber,
			OfferedItemID:        swapRequest.OfferedItemID.String(),
			RequestedItemID:      swapRequest.RequestedItemID.String(),
			OfferedItemOwnerID:   swapRequest.OfferedItemOwnerID.String(),
			RequestedItemOwnerID: swapRequest.RequestedItemOwnerID.String(),
		},
	}

	context.JSON(status, response)
}

func respondWithSwapRequestList(context *gin.Context, status int, message string, swapRequests []domain.SwapRequest) {
	var responseList []SwapRequestResponse

	for _, request := range swapRequests {
		responseList = append(responseList, SwapRequestResponse{
			ID:                   request.ID.String(),
			Status:               string(request.Status),
			ReferenceNumber:      request.ReferenceNumber,
			OfferedItemID:        request.OfferedItemID.String(),
			RequestedItemID:      request.RequestedItemID.String(),
			OfferedItemOwnerID:   request.OfferedItemOwnerID.String(),
			RequestedItemOwnerID: request.RequestedItemOwnerID.String(),
		})
	}

	context.JSON(status, SwapRequestListResponse{
		Message:      message,
		SwapRequests: responseList,
	})
}

func generateReferenceNumber() (string, error) {
	bytes := make([]byte, 12)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
