package influxdb

import "github.com/labstack/echo"

type ProcessApiRestEndp struct {
	Processes []Process
	Echo      *echo.Echo
}

func (endp *ProcessApiRestEndp) SetupRouts() {
	endp.Echo.PUT("/blackflowint/influxdb/api/proc", endp.AddProcessEndpoint)
	endp.Echo.POST("/blackflowint/influxdb/api/proc/ctl", endp.CtlProcessEndpoint)
	endp.Echo.GET("/blackflowint/influxdb/api/proc/filters", endp.GetFiltersEndpoint)
	endp.Echo.PUT("/blackflowint/influxdb/api/proc/filters", endp.AddFilterEndpoint)
	endp.Echo.DELETE("/blackflowint/influxdb/api/proc/filters/:id", endp.RemoveFilterEndpoint)
	endp.Echo.GET("/blackflowint/influxdb/api/proc/selectors", endp.GetSelectorsEndpoint)
	endp.Echo.PUT("/blackflowint/influxdb/api/proc/selectors", endp.AddSelectorEndpoint)
	endp.Echo.DELETE("/blackflowint/influxdb/api/proc/selectors/:id", endp.RemoveSelectorEndpoint)
}

func (endp *ProcessApiRestEndp) AddFilterEndpoint(c echo.Context) error {
	return nil
}
func (endp *ProcessApiRestEndp) RemoveFilterEndpoint(c echo.Context) error {
	return nil
}
func (endp *ProcessApiRestEndp) GetFiltersEndpoint(c echo.Context) error {
	// resp := []Filter{}
	// return c.JSON(http.StatusCreated, endp.Process.GetFilters())
	return nil
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
func (endp *ProcessApiRestEndp) CtlProcessEndpoint(c echo.Context) error {
	return nil
}
func (endp *ProcessApiRestEndp) AddProcessEndpoint(c echo.Context) error {
	return nil
}
func (endp *ProcessApiRestEndp) RemoveProcessEndpoint(c echo.Context) error {
	return nil
}
