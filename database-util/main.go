package main

import (
	"flag"
	"github.com/ghthor/database/config"
	"log"
	"os"
	"os/exec"
)

func main() {
	configFilepath := flag.String("config", "config.json", "Path to a database configuration file")
	requirePwd := flag.Bool("require-password", false, "require the password to be typed to stdin")

	flag.Parse()

	cfg, err := config.ReadFromFile(*configFilepath)
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	var cmd *exec.Cmd
	if *requirePwd {
		cmd = exec.Command("mysqldump", "-d", "-u", cfg.Username, "-p", cfg.DefaultDB)
	} else {
		cmd = exec.Command("mysqldump", "-d", "-u", cfg.Username, "-p"+cfg.Password, cfg.DefaultDB)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
