package influxdb

// AddFilter adds new filter entry
func (pr *Process) AddFilter(filter Filter) IDt {
	defer func() {
		pr.apiMutex.Unlock()
	}()
	pr.apiMutex.Lock()
	filter.ID = GetNewID(pr.Config.Filters)
	pr.Config.Filters = append(pr.Config.Filters, filter)
	return filter.ID
}

// RemoveFilter removes 1 filter entry by ID
func (pr *Process) RemoveFilter(ID IDt) {
	defer func() {
		pr.apiMutex.Unlock()
	}()
	pr.apiMutex.Lock()
	for i, f := range pr.Config.Filters {
		if f.ID == ID {
			pr.Config.Filters = append(pr.Config.Filters[:i], pr.Config.Filters[i+1:]...)
		}
	}
}

// AddSelector adds new selector
func (pr *Process) AddSelector(selector Selector) {
	defer func() {
		pr.apiMutex.Unlock()
	}()
	pr.apiMutex.Lock()
	pr.Config.Selectors = append(pr.Config.Selectors, selector)
	pr.mqttAdapter.Subscribe(selector.Topic, 0)
}

// RemoveSelector removes 1 selector entry
func (pr *Process) RemoveSelector(ID IDt) {
	defer func() {
		pr.apiMutex.Unlock()
	}()
	pr.apiMutex.Lock()
	for i, s := range pr.Config.Selectors {
		if s.ID == ID {
			pr.mqttAdapter.Unsubscribe(s.Topic)
			pr.Config.Selectors = append(pr.Config.Selectors[:i], pr.Config.Selectors[i+1:]...)
		}
	}
}

// GetFilters returns all filters
func (pr *Process) GetFilters() []Filter {
	return pr.Config.Filters
}

// GetSelectors return all selectors
func (pr *Process) GetSelectors() []Selector {
	return pr.Config.Selectors
}
