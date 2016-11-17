package influxdb

import "testing"

func TestIntegration_LoadIntConfigs(t *testing.T) {
	type fields struct {
		Processes      []Process
		ProcessConfigs ProcessConfigs
		StoreLocation  string
		storeFullPath  string
		Name           string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		it := &Integration{
			Processes:      tt.fields.Processes,
			ProcessConfigs: tt.fields.ProcessConfigs,
			StoreLocation:  tt.fields.StoreLocation,
			storeFullPath:  tt.fields.storeFullPath,
			Name:           tt.fields.Name,
		}
		if err := it.LoadConfig(); (err != nil) != tt.wantErr {
			t.Errorf("%q. Integration.LoadIntConfigs() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestIntegration_SaveConfigs(t *testing.T) {

	it := &Integration{
		ProcessConfigs: ProcessConfigs{ProcessConfig{
			Filters:   []Filter{Filter{}},
			Selectors: []Selector{Selector{}},
		},
			ProcessConfig{}},
		StoreLocation: "./",
	}
	it.boot()
	err := it.SaveConfigs()
	if err != nil {
		t.Error(err)
	}

}

func TestIntegration_boot(t *testing.T) {
	type fields struct {
		Processes      []Process
		ProcessConfigs ProcessConfigs
		StoreLocation  string
		storeFullPath  string
		Name           string
	}
	tests := []struct {
		name   string
		fields fields
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		it := &Integration{
			Processes:      tt.fields.Processes,
			ProcessConfigs: tt.fields.ProcessConfigs,
			StoreLocation:  tt.fields.StoreLocation,
			storeFullPath:  tt.fields.storeFullPath,
			Name:           tt.fields.Name,
		}
		it.boot()
	}
}
