package helpers

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var Config struct {
	ParseInterval string `envconfig:"PARSE_INTERVAL" default:"@every 1m"`
	BackupDir string `envconfig:"BACKUP_DIR" default:"./data"`
	Database struct {
		Dialect string `envconfig:"DB_DIALECT" default:"sqlite3"`
		Uri string `envconfig:"DB_URI" default:":memory:"`
	}
}

func init() {
	if err := envconfig.Process("velobike", &Config); err != nil {
		log.Fatal("cannot parse config:", err)
	}
}
