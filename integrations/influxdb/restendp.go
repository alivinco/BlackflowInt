package influxdb

import (
	"net/http"

	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type DefaultResponse struct {
	ID   IDt
	Code int
	Msg  string
}

type ProcCtlRequest struct {
	Action string
}

type IntegrationAPIRestEndp struct {
	integr *Integration
	Echo   *echo.Echo
}

func (endp *IntegrationAPIRestEndp) SetupRoutes() {
	endp.Echo.GET("/blackflowint/influxdb/api/proc/:id", endp.getProcessEndpoint)
	endp.Echo.PUT("/blackflowint/influxdb/api/proc", endp.addProcessEndpoint)
	endp.Echo.POST("/blackflowint/influxdb/api/proc/:id/ctl", endp.ctlProcessEndpoint)
	endp.Echo.GET("/blackflowint/influxdb/api/proc/:id/filters", endp.getFiltersEndpoint)
	endp.Echo.PUT("/blackflowint/influxdb/api/proc/:id/filters", endp.addFilterEndpoint)
	endp.Echo.DELETE("/blackflowint/influxdb/api/proc/:id/filters/:fid", endp.removeFilterEndpoint)
	endp.Echo.GET("/blackflowint/influxdb/api/proc/:id/selectors", endp.getSelectorsEndpoint)
	endp.Echo.PUT("/blackflowint/influxdb/api/proc/:id/selectors", endp.addSelectorEndpoint)
	endp.Echo.DELETE("/blackflowint/influxdb/api/proc/:id/selectors/:sid", endp.removeSelectorEndpoint)
}
func (endp *IntegrationAPIRestEndp) ctlProcessEndpoint(c echo.Context) error {
	req := ProcCtlRequest{}
	procID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	if err = c.Bind(&req); err != nil {
		return err
	}
	switch req.Action {
	case "start":
		err = endp.integr.GetProcessByID(IDt(procID)).Start()
		break
	case "stop":
		err = endp.integr.GetProcessByID(IDt(procID)).Stop()
		break
	case "broker_auto_config":
		endp.integr.BrokerAutoConfig(IDt(procID))
		break
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, DefaultResponse{ID: IDt(procID), Msg: "Action executed"})
}

func (endp *IntegrationAPIRestEndp) addFilterEndpoint(c echo.Context) error {
	filter := Filter{}
	procID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	if err := c.Bind(&filter); err != nil {
		return err
	}
	proc := endp.integr.GetProcessByID(IDt(procID))
	newID := proc.AddFilter(filter)
	endp.integr.SaveConfigs()
	return c.JSON(http.StatusOK, DefaultResponse{ID: newID, Msg: "Filter added."})
}
func (endp *IntegrationAPIRestEndp) removeFilterEndpoint(c echo.Context) error {
	procID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	filterID, err := strconv.Atoi(c.Param("fid"))
	if err != nil {
		return err
	}
	proc := endp.integr.GetProcessByID(IDt(procID))
	proc.RemoveFilter(IDt(filterID))
	endp.integr.SaveConfigs()
	return c.JSON(http.StatusOK, DefaultResponse{ID: IDt(procID), Msg: "Filter removed."})
}
func (endp *IntegrationAPIRestEndp) getProcessEndpoint(c echo.Context) error {
	resp := []ProcessConfig{}
	// var procID IDt
	// var err error
	procID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusNoContent)
	}
	log.Infof("Getting filters for process = %d", IDt(procID))
	for _, proc := range endp.integr.processConfigs {
		if (proc.ID == IDt(procID)) || (proc.ID == 0) {
			// log.Debugf("Proc = %+v", proc.Config)
			resp = append(resp, proc)
		}

	}
	return c.JSON(http.StatusCreated, resp)
}
func (endp *IntegrationAPIRestEndp) getFiltersEndpoint(c echo.Context) error {
	resp := []Filter{}
	// var procID IDt
	// var err error
	procID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusNoContent)
	}
	log.Infof("Getting filters for process = %d", IDt(procID))
	for _, proc := range endp.integr.processes {
		if (proc.Config.ID == IDt(procID)) || (proc.Config.ID == 0) {
			// log.Debugf("Proc = %+v", proc.Config)
			resp = append(resp, proc.GetFilters()...)
		}

	}
	return c.JSON(http.StatusCreated, resp)
	// return nil
}
func (endp *IntegrationAPIRestEndp) getSelectorsEndpoint(c echo.Context) error {
	resp := []Selector{}
	// var procID IDt
	// var err error
	procID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusNoContent)
	}
	log.Infof("Getting filters for process = %d", IDt(procID))
	for _, proc := range endp.integr.processes {
		if (proc.Config.ID == IDt(procID)) || (proc.Config.ID == 0) {
			// log.Debugf("Proc = %+v", proc.Config)
			resp = append(resp, proc.GetSelectors()...)
		}

	}
	return c.JSON(http.StatusCreated, resp)
}
func (endp *IntegrationAPIRestEndp) addSelectorEndpoint(c echo.Context) error {
	selector := Selector{}
	procID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	if err := c.Bind(&selector); err != nil {
		return err
	}
	proc := endp.integr.GetProcessByID(IDt(procID))
	newID := proc.AddSelector(selector)
	endp.integr.SaveConfigs()
	return c.JSON(http.StatusOK, DefaultResponse{ID: newID, Msg: "Selector added."})
}
func (endp *IntegrationAPIRestEndp) removeSelectorEndpoint(c echo.Context) error {
	procID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	selectorID, err := strconv.Atoi(c.Param("sid"))
	if err != nil {
		return err
	}
	proc := endp.integr.GetProcessByID(IDt(procID))
	proc.RemoveSelector(IDt(selectorID))
	endp.integr.SaveConfigs()
	return c.JSON(http.StatusOK, DefaultResponse{ID: IDt(procID), Msg: "Selector removed."})
}

func (endp *IntegrationAPIRestEndp) addProcessEndpoint(c echo.Context) error {
	return nil
}
func (endp *IntegrationAPIRestEndp) removeProcessEndpoint(c echo.Context) error {
	return nil
}
