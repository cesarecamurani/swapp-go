package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"swapp-go/cmd/internal/adapters/handlers/responses"
	"swapp-go/cmd/internal/application/services"
	"swapp-go/cmd/internal/domain"
)

type ItemHandler struct {
	itemService services.ItemServiceInterface
}

func NewItemHandler(itemServiceInterface services.ItemServiceInterface) *ItemHandler {
	return &ItemHandler{itemServiceInterface}
}

type CreateItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PictureURL  string `json:"picture"`
}

type UpdateItemRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	PictureURL  string `json:"picture,omitempty"`
}

type ItemResponse struct {
	ItemID      string `json:"item_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PictureURL  string `json:"picture"`
	UserID      string `json:"user_id"`
}

type ItemSuccessResponse struct {
	Message string        `json:"message"`
	Item    *ItemResponse `json:"item"`
}

func (itemHandler *ItemHandler) CreateItem(context *gin.Context) {
	userID := context.GetString("userID")
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		responses.BadRequest(context, "Invalid user ID", err)
		return
	}

	name := context.PostForm("name")
	description := context.PostForm("description")

	pictureURL, err := saveUploadedPicture(context, "picture")
	if err != nil {
		responses.InternalServerError(context, "Failed to save uploaded picture", err)
	}

	item := &domain.Item{
		Name:        name,
		Description: description,
		PictureURL:  pictureURL,
		UserID:      parsedUserID,
	}

	if err = itemHandler.itemService.CreateItem(item); err != nil {
		responses.BadRequest(context, "Item creation failed", err)
		return
	}

	respondWithItem(context, http.StatusCreated, "Item created successfully!", item)
}

func (itemHandler *ItemHandler) UpdateItem(context *gin.Context) {
	item, ok := itemHandler.verifyItemOwnership(context)
	if !ok {
		return
	}

	name := context.PostForm("name")
	description := context.PostForm("description")

	updateData := make(map[string]interface{})
	if name != "" {
		updateData["name"] = name
	}
	if description != "" {
		updateData["description"] = description
	}

	if url, err := saveUploadedPicture(context, "picture"); err == nil {
		updateData["picture"] = url
	}

	if len(updateData) == 0 {
		responses.BadRequest(context, "No valid fields provided for update", nil)
		return
	}

	updatedItem, err := itemHandler.itemService.UpdateItem(item.ID, updateData)
	if err != nil {
		responses.InternalServerError(context, "Failed to update item", err)
		return
	}

	respondWithItem(context, http.StatusOK, "Item updated successfully!", updatedItem)
}

func (itemHandler *ItemHandler) DeleteItem(context *gin.Context) {
	item, ok := itemHandler.verifyItemOwnership(context)
	if !ok {
		return
	}

	if err := itemHandler.itemService.DeleteItem(item.ID); err != nil {
		responses.BadRequest(context, "Failed to delete item", err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully!"})
}

func (itemHandler *ItemHandler) GetItemByID(context *gin.Context) {
	itemID := context.Param("id")

	parsedID, err := uuid.Parse(itemID)
	if err != nil {
		responses.BadRequest(context, "Invalid item ID", err)
		return
	}

	item, err := itemHandler.itemService.GetItemByID(parsedID)
	if err != nil {
		responses.NotFound(context, "Item not found", err)
		return
	}

	response := &ItemResponse{
		ItemID:      item.ID.String(),
		Name:        item.Name,
		Description: item.Description,
		PictureURL:  item.PictureURL,
		UserID:      item.UserID.String(),
	}

	context.JSON(http.StatusOK, response)
}

func respondWithItem(context *gin.Context, status int, message string, item *domain.Item) {
	response := ItemSuccessResponse{
		Message: message,
		Item: &ItemResponse{
			ItemID:      item.ID.String(),
			Name:        item.Name,
			Description: item.Description,
			PictureURL:  item.PictureURL,
			UserID:      item.UserID.String(),
		},
	}

	context.JSON(status, response)
}

func (itemHandler *ItemHandler) verifyItemOwnership(context *gin.Context) (*domain.Item, bool) {
	itemID := context.Param("id")
	userID := context.GetString("userID")

	parsedID, err := uuid.Parse(itemID)
	if err != nil {
		responses.BadRequest(context, "Invalid item ID", err)
		return nil, false
	}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		responses.BadRequest(context, "Invalid user ID", err)
		return nil, false
	}

	item, err := itemHandler.itemService.GetItemByID(parsedID)
	if err != nil {
		responses.NotFound(context, "Item not found", err)
		return nil, false
	}
	if item.UserID != parsedUserID {
		responses.Unauthorized(context, "This item doesn't belong to you!", nil)
		return nil, false
	}

	return item, true
}

func saveUploadedPicture(context *gin.Context, formKey string) (string, error) {
	file, err := context.FormFile(formKey)
	if err != nil {
		return "", err
	}

	filename := uuid.New().String() + "_" + filepath.Base(file.Filename)
	savePath := filepath.Join("uploads", filename)
	if err = context.SaveUploadedFile(file, savePath); err != nil {
		return "", err
	}

	return "/uploads/" + filename, nil
}
