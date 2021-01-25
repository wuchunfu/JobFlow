package configModel

//type YamlSetting struct {
//	System   System   `mapstructure:"system" json:"system" yaml:"system"`
//	Database Database `mapstructure:"database" json:"database" yaml:"database"`
//	Log      Log      `mapstructure:"log" json:"log" yaml:"log"`
//}
//
//type System struct {
//	AppName          string `mapstructure:"appName" json:"appName" yaml:"appName"`
//	HttpAddr         string `mapstructure:"httpAddr" json:"httpAddr" yaml:"httpAddr"`
//	HttpPort         string `mapstructure:"httpPort" json:"httpPort" yaml:"httpPort"`
//	AllowIps         string `mapstructure:"allowIps" json:"allowIps" yaml:"allowIps"`
//	ConcurrencyQueue int    `mapstructure:"concurrencyQueue" json:"concurrencyQueue" yaml:"concurrencyQueue"`
//	ApiSignEnable    bool   `mapstructure:"apiSignEnable" json:"apiSignEnable" yaml:"apiSignEnable"`
//	ApiKey           string `mapstructure:"apiKey" json:"apiKey" yaml:"apiKey"`
//	ApiSecret        string `mapstructure:"apiSecret" json:"apiSecret" yaml:"apiSecret"`
//	AuthSecret       string `mapstructure:"authSecret" json:"authSecret" yaml:"authSecret"`
//	EnableTls        bool   `mapstructure:"enableTls" json:"enableTls" yaml:"enableTls"`
//	CaFile           string `mapstructure:"caFile" json:"caFile" yaml:"caFile"`
//	CertFile         string `mapstructure:"certFile" json:"certFile" yaml:"certFile"`
//	KeyFile          string `mapstructure:"keyFile" json:"keyFile" yaml:"keyFile"`
//}
//
//type Database struct {
//	DbType       string `mapstructure:"dbType" json:"dbType" yaml:"dbType"`
//	Host         string `mapstructure:"host" json:"host" yaml:"host"`
//	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
//	DbName       string `mapstructure:"dbName" json:"dbName" yaml:"dbName"`
//	Username     string `mapstructure:"username" json:"username" yaml:"username"`
//	Password     string `mapstructure:"password" json:"password" yaml:"password"`
//	Prefix       string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
//	Charset      string `mapstructure:"charset" json:"charset" yaml:"charset"`
//	MaxIdleConns int    `mapstructure:"maxIdleConns" json:"maxIdleConns" yaml:"maxIdleConns"`
//	MaxOpenConns int    `mapstructure:"maxOpenConns" json:"maxOpenConns" yaml:"maxOpenConns"`
//	Level        string `mapstructure:"level" json:"level" yaml:"level"`
//	SslMode      string `mapstructure:"sslMode" json:"sslMode" yaml:"sslMode"`
//	TimeZone     string `mapstructure:"timeZone" json:"timeZone" yaml:"timeZone"`
//}
//
//type Log struct {
//	FilePath string `mapstructure:"filePath" json:"path" yaml:"filePath"`
//	FileName string `mapstructure:"fileName" json:"fileName" yaml:"fileName"`
//	Level    string `mapstructure:"level" json:"level" yaml:"level"`
//	Mode     string `mapstructure:"mode" json:"mode" yaml:"mode"`
//}
