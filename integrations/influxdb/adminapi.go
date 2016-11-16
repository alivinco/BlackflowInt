package influxdb

// AddFilter adds new filter entry
func (pr *Process) AddFilter(filter Filter) int {
	defer func() {
		pr.apiMutex.Unlock()
	}()
	pr.apiMutex.Lock()
	var newID int
	for _, f := range pr.filters {
		if f.ID > newID {
			newID = f.ID
		}
	}

	newID++
	filter.ID = newID
	pr.filters = append(pr.filters, filter)
	return newID
}

// RemoveFilter removes 1 filter entry by ID
func (pr *Process) RemoveFilter(ID int) {
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
func (pr *Process) RemoveSelector(selector *Selector) {
	defer func() {
		pr.apiMutex.Unlock()
	}()
	pr.apiMutex.Lock()
	pr.mqttAdapter.Unsubscribe(selector.Topic)
	for i, s := range pr.selectors {
		if s.Topic == selector.Topic {
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
