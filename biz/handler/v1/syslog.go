package v1

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"net/http"
	"os"
	"sync"

	"github.com/colin-404/logx"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	syslogWriter *syslog.Writer
	once         sync.Once
	syslogErr    error
)

type SyslogConfig struct {
	Host     string
	Port     string
	Tag      string
	Protocol string
	Level    string
	User     string
	Password string
}

func InitSyslog() {

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

	syslogWriter, syslogErr = syslog.Dial(Protocol, fmt.Sprintf("%s:%s", Host, Port), syslogLevel, "webhook")
	if syslogErr != nil {
		logx.Errorf("Failed to dial syslog: %v", syslogErr)
		os.Exit(1)
	}
}

func Webhook(c *gin.Context) {
	once.Do(InitSyslog) // Ensure syslog is initialized only once

	rawData, err := c.GetRawData()
	if err != nil {
		logx.Errorf("Error reading request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read request body"})
		return
	}

	// Log the raw JSON data to your file log (wazuh-webhook.log)
	logx.Infof("Webhook received raw data: %s", string(rawData))

	// Attempt to parse the JSON to validate it and potentially extract specific fields
	var payload map[string]interface{} // Use a map for flexible JSON structure
	if err := json.Unmarshal(rawData, &payload); err != nil {
		logx.Errorf("Error unmarshalling JSON: %v. Raw data: %s", err, string(rawData))
		// Decide if you still want to send raw data to syslog if unmarshalling fails
	}

	infoData := fmt.Sprintf("Received JSON data: event-type: %s, event_data: %s", payload["event_type"], string(rawData))

	logx.Infof("infoData: %s", infoData)

	// Send the raw data to Syslog
	if syslogWriter != nil {

		err = syslogWriter.Info(infoData)
		if err != nil {
			logx.Errorf("Error writing to syslog: %v", err)
			// Potentially retry or handle syslog unavailability
		}
	} else {
		logx.Errorf("Syslog writer is not initialized. Cannot send log.")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"status":  "received",
	})
}
