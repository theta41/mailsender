package cfg

type Config struct {
	Debug           bool `yaml:"debug"`
	PeriodMinutes   int  `yaml:"period_minutes"`
	RequstsInPeriod int  `yaml:"requests"`
	Kafka           struct {
		Brokers []string `yaml:"brokers"`
		Topic   string   `yaml:"topic"`
		GroupId string   `yaml:"group_id"`
	} `yaml:"kafka"`
}

func NewConfig(yamlFile string) (conf Config, err error) {
	err = loadYaml(yamlFile, &conf)
	return
}
