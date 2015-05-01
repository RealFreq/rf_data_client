package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func init() {
	viper.SetConfigName("config")
	viper.ReadInConfig()
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Scan frequencies from 30MHz to 1GHz, use 100Hz bin size for FFT
	// Integration interval is 5 minutes
	// Gain is 50
	cmd := exec.Command("rtl_power", "-f", "30M:1.1G:100", "-i", "5s", "-g", "50", "-1")

	go func() {
		sig := <-sigs
		log.Printf("Caught signal (%s), exiting!\n", sig)
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill: ", err)
		}
	}()

	log.Print("Collecting RF data...")

	result, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed\nError running command: %s\n", err)
	}

	log.Println("finished")

	log.Print("Connecting to server...")

	config := viper.GetStringMap("rf_server")
	server := fmt.Sprintf("%s:%d", config["host"].(string), config["port"].(int))

	conn, err := net.Dial("tcp", server)
	if err != nil {
		// TODO save data for later upload
		log.Printf("failed\nError connecting to server: %s\n", err)
	}

	log.Print("connected")

	log.Println("Starting upload...")

	for _, line := range strings.Split(string(result[:]), "\n") {
		fmt.Fprintf(conn, line+"\n")
	}

	log.Println("Upload finished")
}
