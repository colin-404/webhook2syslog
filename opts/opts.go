package opts

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/colin-404/logx"
	"github.com/spf13/viper"
)

func initSyslog() *syslog.Writer {

	//如果Protocol 、Host、Port为空，则报错并关闭程序

	Host := viper.GetString("syslog.host")
	Port := viper.GetString("syslog.port")
	Protocol := viper.GetString("syslog.protocol")
	Level := viper.GetString("syslog.level")

	if Host == "" || Port == "" || Protocol == "" || Level == "" {
		logx.Errorf("syslog.host, syslog.port, syslog.protocol, syslog.level is required")
		os.Exit(1)
	}

	syslogLevel := syslog.LOG_INFO
	if Level == "info" {
		syslogLevel = syslog.LOG_INFO
	} else if Level == "debug" {
		syslogLevel = syslog.LOG_DEBUG
	} else if Level == "error" {
		syslogLevel = syslog.LOG_ERR
	} else if Level == "warning" {
		syslogLevel = syslog.LOG_WARNING
	} else {
		logx.Errorf("syslog.level is invalid")
		os.Exit(1)
	}

	//获取syslog配置
	// syslogConfig := SyslogConfig{
	// 	Host:     viper.GetString("syslog.host"),
	// 	Port:     viper.GetString("syslog.port"),
	// 	Protocol: viper.GetString("syslog.protocol"),
	// }
	// syslog.Dial()

	syslogWriter, syslogErr := syslog.Dial(Protocol, fmt.Sprintf("%s:%s", Host, Port), syslogLevel, "webhook")
	if syslogErr != nil {
		logx.Errorf("Failed to dial syslog: %v", syslogErr)
		// Handle error appropriately - maybe the app shouldn't start or should retry
	}

	return syslogWriter
}
