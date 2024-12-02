package lokinet

import (
	"fmt"
	"os/exec"
)

func StartLokinet() {
	cmd := exec.Command("lokinet", "start")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to start Lokinet:", err)
	} else {
		fmt.Println("Lokinet started successfully!")
	}
}
