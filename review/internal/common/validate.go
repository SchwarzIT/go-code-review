package common

import (
	"log"
	"os"
	"runtime"
)

func ValidateCPUs() {
	if os.Getenv("SKIP_CPU_CHECK") == "1" {
		return
	}
	log.Fatalf("%s", os.Getenv("SKIP_CPU_CHECK"))
	if runtime.NumCPU() != 32 {
		log.Fatal("This api is meant to be run on 32 core machines. Exiting...")
	}
}
