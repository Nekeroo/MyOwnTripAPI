package main

import (
	"fmt"
	"os"
)

func main() {

	cfg := config{
		addr: ":3000",
	}

	app := application{
		config: cfg,
	}

	if err := app.run(app.mount()); err != nil {
		fmt.Printf("Server can't start, err : %s", err)
		os.Exit(1)
	}

}
