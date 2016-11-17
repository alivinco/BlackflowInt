package influxdb

import "github.com/labstack/echo"
import "net/http"

type ProcessApiRestEndp struct {
	Process *Process
	Echo    *echo.Echo
}

func (endp *ProcessApiRestEndp) SetupRouts() {
	endp.Echo.GET("/blackflowint/influxdb/api/filters", endp.GetFiltersEndpoint)
	endp.Echo.GET("/blackflowint/influxdb/api/selectors", endp.GetSelectorsEndpoint)
	endp.Echo.PUT("/blackflowint/influxdb/api/selectors", endp.AddSelectorEndpoint)
	endp.Echo.DELETE("/blackflowint/influxdb/api/selectors", endp.RemoveSelectorEndpoint)
	endp.Echo.PUT("/blackflowint/influxdb/api/filters", endp.AddFilterEndpoint)
	endp.Echo.DELETE("/blackflowint/influxdb/api/filters/:id", endp.RemoveFilterEndpoint)
}

func (endp *ProcessApiRestEndp) AddFilterEndpoint(c echo.Context) error {
	return nil
}
func (endp *ProcessApiRestEndp) RemoveFilterEndpoint(c echo.Context) error {
	return nil
}
func (endp *ProcessApiRestEndp) GetFiltersEndpoint(c echo.Context) error {
	return c.JSON(http.StatusCreated, endp.Process.GetFilters())
}
func (endp *ProcessApiRestEndp) GetSelectorsEndpoint(c echo.Context) error {
	return nil
}
func (endp *ProcessApiRestEndp) AddSelectorEndpoint(c echo.Context) error {
	return nil
}
func (endp *ProcessApiRestEndp) RemoveSelectorEndpoint(c echo.Context) error {
	return nil
}
