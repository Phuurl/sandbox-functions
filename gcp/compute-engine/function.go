// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"log"
	"os"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func FunctionEntry(ctx context.Context, m PubSubMessage) error {

	oauthClient, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(oauthClient)
	if err != nil {
		log.Fatalf("Unable to create Compute service: %v\n", err)
	}

	// Pull in the current project ID from the Cloud Function environment variables
	projectID := os.Getenv("GCLOUD_PROJECT")
	stopCount := 0

	zoneList := computeService.Zones.List(projectID)
	if zoneErr := zoneList.Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			instanceList := computeService.Instances.List(projectID, zone.Name)
			if err := instanceList.Pages(ctx, func(page *compute.InstanceList) error {
				for _, instance := range page.Items {
					if instance.Status == "RUNNING" {
						_, preserve := instance.Labels["preserve"] // Check for 'preserve' key (regardless of value)
						if !preserve {
							fmt.Printf("Stopping running instance %s in %s\n", instance.Name, zone.Name)
							_, err := computeService.Instances.Stop(projectID, zone.Name, instance.Name).Do()
							if err != nil {
								log.Fatal(err)
							} else {
								stopCount++
							}
						}
					} else {
						fmt.Printf("Non-running instance %s detected in %s\n", instance.Name, zone.Name)
					}
				}
				return nil
			}); err != nil {
				log.Fatalf("Error in %s zone: %s\n", zone.Name, err)
			}
		}
		return nil
	}); zoneErr != nil {
		log.Fatal(zoneErr)
	}
	fmt.Printf("Stopped %d VMs\n", stopCount)
	return nil
}
