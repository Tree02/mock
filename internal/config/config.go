package config

import (
	"fmt"
	domainPuesto "mockLogin/internal/domain/puesto"
	domainServer "mockLogin/internal/domain/server"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var GlobalConfig *Config
var once sync.Once

type Config struct {
	Server  domainServer.ServerConfig   `mapstructure:"server"` // Configuración de server cargada desde archivo de configuración
	Puestos map[int]domainPuesto.Puesto `mapstructure:"-"`      // Puestos cargados desde archivo de configuración
}

// se crea estructura interna con el fin de agregar al mapa de puestos
type RawPuesto struct {
	Id           int    `mapstructure:"id"`
	AgencyBranch string `mapstructure:"agencyBranch"`
	AgencyCode   string `mapstructure:"agencyCode"`
	Workstation  string `mapstructure:"workstation"`
}

func Load() {
	once.Do(func() {
		// setea configuraciones de viper
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")

		// lee archivo
		configPath := fileExtend()
		if configPath == "" {
			fmt.Println("No se encontró archivo de configuración compatible.")
			return
		}
		ext := strings.TrimPrefix(filepath.Ext(configPath), ".")
		viper.SetConfigType(ext)

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error al leer archivo\n")
		}

		// manejo de estructura interna mapeando puestos
		var rawConfig struct {
			Server  domainServer.ServerConfig `mapstructure:"server"`
			Puestos []RawPuesto               `mapstructure:"puestos"`
		}
		if err := viper.Unmarshal(&rawConfig); err != nil {
			fmt.Println("Error al mapear a config")
			return
		}

		puestosMap := make(map[int]domainPuesto.Puesto)
		for _, rawPuesto := range rawConfig.Puestos {
			puestosMap[rawPuesto.Id] = domainPuesto.Puesto{
				AgencyBranch: rawPuesto.AgencyBranch,
				AgencyCode:   rawPuesto.AgencyCode,
				Workstation:  rawPuesto.Workstation,
			}
		}

		// asigna valores a variable para ser retornada y manejada a nivel global
		GlobalConfig = &Config{
			Server:  rawConfig.Server,
			Puestos: puestosMap,
		}
		fmt.Println("Cargado todo bien desde viper")
	})
}

func GetConfig() *Config {
	return GlobalConfig
}

// detecta extensión de archivo config.
func fileExtend() string {
	opts := []string{
		"config.yaml",
		"config.yam",
		"config.json",
	}

	for _, candidate := range opts {
		if fileExists(candidate) {
			return candidate
		}
	}
	return ""
}

// busca archivo
func fileExists(filename string) bool {
	_, err := filepath.Abs(filename)
	return err == nil
}
