package geo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rodinv/errors"

	"github.com/rodinv/geodata/internal/geo/repository"
)

type Handler struct {
	repo RepositoryProvider
}

type RepositoryProvider interface {
	GetLocationByIP(ip string) (*repository.GeoItem, error)
	GetLocationsByCity(city string) ([]*repository.GeoItem, error)
}

func New(repo RepositoryProvider) *Handler {
	return &Handler{
		repo: repo,
	}
}

// GetLocationByIP godoc
// @Summary      Get location by IP
// @Description  get location by IP
// @Produce      json
// @Param        ip   query      string  true  "ip"
// @Success      200  {object}  repository.GeoItem
// @Failure      400  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /ip/location [get]
func (h *Handler) GetLocationByIP(ctx echo.Context) error {
	// check input
	ip := ctx.QueryParam("ip")
	if ip == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ip param is empty")
	}

	// get data
	item, err := h.repo.GetLocationByIP(ip)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	case err != nil:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, item)
}

// GetLocations godoc
// @Summary      Get locations by city
// @Description  get location by city
// @Produce      json
// @Param        city   query   string  true  "city"
// @Success      200  	{array}   repository.GeoItem
// @Failure      400  	{object}  echo.HTTPError
// @Failure      404  	{object}  echo.HTTPError
// @Failure      500  	{object}  echo.HTTPError
// @Router       /city/locations [get]
func (h *Handler) GetLocations(ctx echo.Context) error {
	// check input
	city := ctx.QueryParam("city")
	if city == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "city param is empty")
	}

	// get data
	items, err := h.repo.GetLocationsByCity(city)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	case err != nil:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, items)
}
