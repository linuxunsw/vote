package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
	"github.com/linuxunsw/vote/tui/internal/tui"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("tui.host", "0.0.0.0")
	viper.SetDefault("tui.port", "2222")
	viper.SetDefault("tui.local", false)
	viper.SetDefault("society", "$ linux society")
	viper.SetDefault("event", "event name")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$XDG_CONFIG_HOME/vote/")
	viper.AddConfigPath("$HOME/.config/vote/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found")
		} else {
			log.Fatalln("Error reading config file: ", err)
		}
	}

	local := viper.GetBool("tui.local")
	host := viper.GetString("tui.host")
	port := viper.GetString("tui.port")

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("error closing log file: %v", err)
		}
	}()

	if !local {
		tui.SSH(host, port)
	} else {
		tui.Local()
	}

}
