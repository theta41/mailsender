package cfg

type Config struct {
	Debug bool `yaml:"debug"`
	Kafka struct {
		User     string   `yaml:"user"`
		Password string   `yaml:"password"`
		Brokers  []string `yaml:"brokers"`
		Topic    string   `yaml:"topic"`
		GroupId  string   `yaml:"group_id"`
	} `yaml:"kafka"`
}

func NewConfig(yamlFile string) (conf Config, err error) {
	err = loadYaml(yamlFile, &conf)
	return
}
