package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/sigdown/kartograf-backend-go/internal/usecase"
)

type handler struct {
	auth   *usecase.AuthService
	points *usecase.PointService
	maps   *usecase.MapService
}

func newHandler(services Services) *handler {
	return &handler{
		auth:   services.Auth,
		points: services.Points,
		maps:   services.Maps,
	}
}

func (h *handler) registerUser(c *gin.Context) {
	var input usecase.RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.auth.Register(c.Request.Context(), input)
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *handler) loginUser(c *gin.Context) {
	var input usecase.LoginUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.auth.Login(c.Request.Context(), input)
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *handler) getCurrentUser(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	user, err := h.auth.Me(c.Request.Context(), claims.UserID)
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *handler) updateAccount(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	var input usecase.UpdateAccountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.auth.UpdateAccount(c.Request.Context(), claims.UserID, input)
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *handler) deleteAccount(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	if err := h.auth.DeleteAccount(c.Request.Context(), claims.UserID); err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *handler) listMaps(c *gin.Context) {
	maps, err := h.maps.List(c.Request.Context())
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, maps)
}

func (h *handler) getMap(c *gin.Context) {
	m, err := h.maps.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, m)
}

func (h *handler) downloadMap(c *gin.Context) {
	mapID, err := parseUUIDParam(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	url, err := h.maps.DownloadURL(c.Request.Context(), mapID)
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *handler) listPoints(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	points, err := h.points.List(c.Request.Context(), claims.UserID)
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, points)
}

func (h *handler) createPoint(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	var input usecase.CreatePointInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	point, err := h.points.Create(c.Request.Context(), claims.UserID, input)
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusCreated, point)
}

func (h *handler) updatePoint(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	pointID, err := parseInt64Param(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var input usecase.UpdatePointInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	point, err := h.points.Update(c.Request.Context(), claims.UserID, pointID, input)
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, point)
}

func (h *handler) deletePoint(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	pointID, err := parseInt64Param(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.points.Delete(c.Request.Context(), claims.UserID, pointID); err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *handler) createMap(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	var input usecase.CreateMapInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	m, err := h.maps.Create(c.Request.Context(), claims.UserID, input)
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusCreated, m)
}

func (h *handler) createMapUploadURL(c *gin.Context) {
	var input usecase.CreateMapUploadInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.maps.StartCreateUpload(c.Request.Context(), input)
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *handler) updateMapMetadata(c *gin.Context) {
	mapID, err := parseUUIDParam(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var input usecase.UpdateMapMetadataInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	m, err := h.maps.UpdateMetadata(c.Request.Context(), mapID, input)
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, m)
}

func (h *handler) replaceMapArchiveUploadURL(c *gin.Context) {
	mapID, err := parseUUIDParam(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var input usecase.ReplaceMapArchiveUploadInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.maps.StartReplaceArchiveUpload(c.Request.Context(), mapID, input)
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *handler) replaceMapArchive(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	mapID, err := parseUUIDParam(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var input usecase.ReplaceMapArchiveInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	archive, err := h.maps.ReplaceArchive(c.Request.Context(), claims.UserID, mapID, input)
	if err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, archive)
}

func (h *handler) deleteMap(c *gin.Context) {
	mapID, err := parseUUIDParam(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.maps.Delete(c.Request.Context(), mapID); err != nil {
		writeError(c, statusFromError(err), err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func parseUUIDParam(value string) (string, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func parseInt64Param(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func parseOptionalInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}
