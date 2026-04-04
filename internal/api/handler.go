package api

import (
	"api-status-check/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	apiService   *service.APIService
	checkService *service.CheckService
}

func NewHandler() *Handler {
	return &Handler{
		apiService:   service.NewAPIService(),
		checkService: service.NewCheckService(),
	}
}

func (h *Handler) CreateAPI(c *gin.Context) {
	var req CreateAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Error(err.Error()))
		return
	}

	config, err := h.apiService.Create(req.Name, req.Tag, req.APIURL, req.Token, req.Model, req.Enabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(config))
}

func (h *Handler) GetAPIs(c *gin.Context) {
	configs, err := h.apiService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(configs))
}

func (h *Handler) GetAPI(c *gin.Context) {
	id := c.Param("id")
	config, err := h.apiService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(config))
}

func (h *Handler) UpdateAPI(c *gin.Context) {
	id := c.Param("id")
	var req UpdateAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Error(err.Error()))
		return
	}

	if err := h.apiService.Update(id, req.Name, req.Tag, req.APIURL, req.Token, req.Model, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(nil))
}

func (h *Handler) DeleteAPI(c *gin.Context) {
	id := c.Param("id")
	if err := h.apiService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(nil))
}

func (h *Handler) GetHistory(c *gin.Context) {
	history, err := h.checkService.GetHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(history))
}

func (h *Handler) GetRecent(c *gin.Context) {
	recent, err := h.checkService.GetRecent100()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(recent))
}

func (h *Handler) RunCheck(c *gin.Context) {
	go h.checkService.RunCheck()
	c.JSON(http.StatusOK, Success("check started"))
}
