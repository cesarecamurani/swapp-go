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

func (handler *ItemHandler) Create(context *gin.Context) {
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

	if err = handler.itemService.Create(item); err != nil {
		responses.BadRequest(context, "Item creation failed", err)
		return
	}

	respondWithItem(context, http.StatusCreated, "Item created successfully!", item)
}

func (handler *ItemHandler) Update(context *gin.Context) {
	item, ok := handler.verifyItemOwnership(context)
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

	updatedItem, err := handler.itemService.Update(item.ID, updateData)
	if err != nil {
		responses.InternalServerError(context, "Failed to update item", err)
		return
	}

	respondWithItem(context, http.StatusOK, "Item updated successfully!", updatedItem)
}

func (handler *ItemHandler) Delete(context *gin.Context) {
	item, ok := handler.verifyItemOwnership(context)
	if !ok {
		return
	}

	if err := handler.itemService.Delete(item.ID); err != nil {
		responses.BadRequest(context, "Failed to delete item", err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully!"})
}

func (handler *ItemHandler) FindByID(context *gin.Context) {
	itemID := context.Param("id")

	parsedID, err := uuid.Parse(itemID)
	if err != nil {
		responses.BadRequest(context, "Invalid item ID", err)
		return
	}

	item, err := handler.itemService.FindByID(parsedID)
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

func (handler *ItemHandler) verifyItemOwnership(context *gin.Context) (*domain.Item, bool) {
	itemID := context.Param("id")
	userID := context.GetString("userID")

	if userID == "" {
		responses.BadRequest(context, "Missing user ID in context", nil)
		return nil, false
	}

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

	item, err := handler.itemService.FindByID(parsedID)
	if err != nil {
		responses.NotFound(context, "Item not found", err)
		return nil, false
	}
	if item == nil {
		responses.InternalServerError(context, "Unexpected nil item returned", nil)
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
