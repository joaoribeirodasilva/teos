package info

import (
	"fmt"
	"strings"
	"time"
)

func Print(servcie string, version string) {

	now := time.Now()
	fmt.Printf("TEOS - Unified Teaching and Learning Platform\n")
	fmt.Printf("%s Micro-Service\n", strings.ToUpper(servcie))
	fmt.Printf("Version %s\n", version)
	fmt.Printf("by BIQX (2024-%s)\n", now.Format("2006"))

}
