package libwebring_test

import (
	"context"
	"fmt"
	"log"

	"libdb.so/libwebring-go"
)

func Example() {
	const webringURL = `https://raw.githubusercontent.com/diamondburned/libwebring/main/example/webring.json`
	ctx := context.Background()

	webring, err := libwebring.FetchData(ctx, webringURL)
	if err != nil {
		log.Fatalln("failed to fetch webring:", err)
	}

	status, err := libwebring.FetchStatusForWebring(ctx, webringURL)
	if err != nil {
		log.Fatalln("failed to fetch status:", err)
	}

	fmt.Print("webring ", webring.Name, " contains these working links:\n")
	for _, link := range webring.Ring.ExcludeAnomalies(status.Anomalies) {
		fmt.Print("  - ", link.Name, " (", link.Link, ")\n")
	}

	// Output:
	// webring acmRing contains these working links:
	//   - diamond (libdb.so)
	//   - aaronlieb (lieber.men)
	//   - etok (etok.codes)
}
