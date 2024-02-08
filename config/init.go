package config

type Init struct {
	Apollo Apollo `mapstructure:"apollo" json:"apollo" yaml:"apollo"`
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis  Redis  `mapstructure:"redis" json:"redis" yaml:"redis"`
}
