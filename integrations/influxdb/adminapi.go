package influxdb

// AddFilter adds new filter entry
func (pr *Process) AddFilter(filter Filter) IDt {
	defer func() {
		pr.apiMutex.Unlock()
	}()
	pr.apiMutex.Lock()
	filter.ID = GetNewID(pr.filters)
	pr.filters = append(pr.filters, filter)
	return filter.ID
}

// RemoveFilter removes 1 filter entry by ID
func (pr *Process) RemoveFilter(ID IDt) {
	defer func() {
		pr.apiMutex.Unlock()
	}()
	pr.apiMutex.Lock()
	for i, f := range pr.filters {
		if f.ID == ID {
			pr.filters = append(pr.filters[:i], pr.filters[i+1:]...)
		}
	}
}

// AddSelector adds new selector
func (pr *Process) AddSelector(selector Selector) {
	defer func() {
		pr.apiMutex.Unlock()
	}()
	pr.apiMutex.Lock()
	pr.selectors = append(pr.selectors, selector)
	pr.mqttAdapter.Subscribe(selector.Topic, 0)
}

// RemoveSelector removes 1 selector entry
func (pr *Process) RemoveSelector(ID IDt) {
	defer func() {
		pr.apiMutex.Unlock()
	}()
	pr.apiMutex.Lock()
	for i, s := range pr.selectors {
		if s.ID == ID {
			pr.mqttAdapter.Unsubscribe(s.Topic)
			pr.selectors = append(pr.selectors[:i], pr.selectors[i+1:]...)
		}
	}
}

// GetFilters returns all filters
func (pr *Process) GetFilters() []Filter {
	return pr.filters
}

// GetSelectors return all selectors
func (pr *Process) GetSelectors() []Filter {
	return pr.filters
}
