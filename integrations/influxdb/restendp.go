package influxdb

import (
	"net/http"

	"github.com/labstack/echo"
)

type IntegrationApiRestEndp struct {
	Processes []*Process
	Echo      *echo.Echo
}

func (endp *IntegrationApiRestEndp) SetupRouts() {
	endp.Echo.PUT("/blackflowint/influxdb/api/proc", endp.AddProcessEndpoint)
	endp.Echo.POST("/blackflowint/influxdb/api/proc/ctl", endp.CtlProcessEndpoint)
	endp.Echo.GET("/blackflowint/influxdb/api/proc/filters", endp.GetFiltersEndpoint)
	endp.Echo.PUT("/blackflowint/influxdb/api/proc/filters", endp.AddFilterEndpoint)
	endp.Echo.DELETE("/blackflowint/influxdb/api/proc/filters/:id", endp.RemoveFilterEndpoint)
	endp.Echo.GET("/blackflowint/influxdb/api/proc/selectors", endp.GetSelectorsEndpoint)
	endp.Echo.PUT("/blackflowint/influxdb/api/proc/selectors", endp.AddSelectorEndpoint)
	endp.Echo.DELETE("/blackflowint/influxdb/api/proc/selectors/:id", endp.RemoveSelectorEndpoint)
}

func (endp *IntegrationApiRestEndp) AddFilterEndpoint(c echo.Context) error {
	return nil
}
func (endp *IntegrationApiRestEndp) RemoveFilterEndpoint(c echo.Context) error {
	return nil
}
func (endp *IntegrationApiRestEndp) GetFiltersEndpoint(c echo.Context) error {
	resp := []Filter{}
	for _, proc := range endp.Processes {
		resp = append(resp, proc.GetFilters()...)
	}
	return c.JSON(http.StatusCreated, resp)
	// return nil
}
func (endp *IntegrationApiRestEndp) GetSelectorsEndpoint(c echo.Context) error {
	return nil
}
func (endp *IntegrationApiRestEndp) AddSelectorEndpoint(c echo.Context) error {
	return nil
}
func (endp *IntegrationApiRestEndp) RemoveSelectorEndpoint(c echo.Context) error {
	return nil
}
func (endp *IntegrationApiRestEndp) CtlProcessEndpoint(c echo.Context) error {
	return nil
}
func (endp *IntegrationApiRestEndp) AddProcessEndpoint(c echo.Context) error {
	return nil
}
func (endp *IntegrationApiRestEndp) RemoveProcessEndpoint(c echo.Context) error {
	return nil
}
