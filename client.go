package main

import (
	//"bufio"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Scan frequencies from 30MHz to 1GHz, use 100Hz bin size for FFT
	// Integration interval is 5 minutes
	// Gain is 50
	// TODO Run this for one integration interval (flag -1), wait an hour,
	//      then run it again
	cmd := exec.Command("rtl_power", "-f", "30M:1.1G:100", "-i", "5s", "-g", "50", "-1")
	//stdout, err := cmd.StdoutPipe()
	//if err != nil {
	//log.Fatal(err)
	//}
	//if err := cmd.Start(); err != nil {
	//log.Fatal(err)
	//}

	go func() {
		sig := <-sigs
		log.Printf("Caught signal (%s), exiting!\n", sig)
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill: ", err)
		}
	}()

	result, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error running command: %s\n", err)
	}

	log.Printf("Result:\n%s\n", result)

	//scanner := bufio.NewScanner(stdout)

	//for scanner.Scan() {
	//line := scanner.Text()
	//log.Println(line)
	//}

	//if err := scanner.Err(); err != nil {
	//log.Printf("Error reading standard input:", err)
	//}
}
