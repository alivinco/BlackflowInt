package models

//MainConfig application root config
type MainConfig struct {
	StorageLocation        string
	AdminRestAPIBindAddres string
	LogLevel               string
	LogPath                string
	ActiveIntegrations     []string
}
