# GCP Compute Engine

Shuts down any running Compute Engine VMs that aren't labeled with `preserve`, across all GCP zones.

## Setup

- Runtime: `Go 1.11`
- Function to execute: `FunctionEntry`
- Memory: `128 MB`

## Trigger

Function is triggered from a Pub/Sub topic. This could then be published to by a Cloud Scheduler running to a cron schedule.

## IAM

This Cloud Function has been tested using a service account with the `Compute Instance Admin (beta)` role
