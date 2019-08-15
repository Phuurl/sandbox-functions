// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"log"
	"os"
	"strings"
	"sync"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

var stopCount = 0
var zoneStopCount = 0
var wg sync.WaitGroup
var debugLogging = false

func DebugLog(message string) {
	if debugLogging {
		fmt.Print(message)
	}
}

func ShutdownZoneVMs(ctx context.Context, computeService *compute.Service, projectID string, zoneID string) {
	instanceList := computeService.Instances.List(projectID, zoneID)
	if err := instanceList.Pages(ctx, func(page *compute.InstanceList) error {
		zoneStop := false
		for _, instance := range page.Items {
			if instance.Status == "RUNNING" {
				// Check for 'preserve' key (regardless of value)
				_, preserve := instance.Labels["preserve"]
				// Stop the VM if not marked for preservation
				if !preserve {
					DebugLog(fmt.Sprintf("Stopping running instance %s in %s\n", instance.Name, zoneID))
					_, err := computeService.Instances.Stop(projectID, zoneID, instance.Name).Do()
					if err != nil {
						log.Fatal(err)
					} else {
						stopCount++
						zoneStop = true
					}
				}
			} else {
				DebugLog(fmt.Sprintf("Non-running instance %s detected in %s\n", instance.Name, zoneID))
			}
		}
		if zoneStop {
			zoneStopCount++
		}
		DebugLog(fmt.Sprintf("goroutine for zone %s ending\n", zoneID))
		return nil
	}); err != nil {
		log.Fatalf("Error in %s zone: %s\n", zoneID, err)
	}
	wg.Done()
}

func FunctionEntry(ctx context.Context, m PubSubMessage) error {
	if strings.EqualFold(os.Getenv("LOG_LEVEL"), "debug") {
		// Log level set to debug
		debugLogging = true
	}

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
	DebugLog(fmt.Sprintf("projectID: %s\n", projectID))

	zoneList := computeService.Zones.List(projectID)
	if zoneErr := zoneList.Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			wg.Add(1)
			go ShutdownZoneVMs(ctx, computeService, projectID, zone.Name)
		}
		return nil
	}); zoneErr != nil {
		log.Fatal(zoneErr)
	}

	// Wait for any goroutines still running
	wg.Wait()
	fmt.Printf("Stopped %d VMs in %d zones\n", stopCount, zoneStopCount)
	return nil
}
