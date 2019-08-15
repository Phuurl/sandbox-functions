# ⛅️ sandbox-functions

A set of FaaS (eg Lambda) functions to enable sandbox behaviours (such as auto stopping resources on a schedule) in your favourite cloud platforms.

## Current support

- AWS (EC2s)
- GCP (Compute Engine VMs)

## Usage

Simply deploy the function as described in the README for the appropriate platform, and configure the trigger to call the function at the desired interval.

As a rule the functions will not touch any resources that are tagged/labelled with `preserve`, to allow exclusions when required.

## Contributing

Pull requests are welcome - please ensure that you update any READMEs where needed.
