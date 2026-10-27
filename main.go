// Command silicondb is Silicon DB: the specifications of processors,
// graphics chips, SoCs, AI accelerators and other data-processing chips.
package main

import (
	"log"

	"github.com/mehran9675/silicondb/internal/ui/app"
)

func main() {
	if err := app.New().Run(); err != nil {
		log.Fatal(err)
	}
}
