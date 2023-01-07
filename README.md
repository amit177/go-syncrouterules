# go-syncrouterules

The project creates routing rules on specific tables for routes based on the next-hop value.

## Use case

I've created this to be able to do network routing based on the source IP.
The table that the routing rules are added to has a static route to a different uplink, hence that allows you to forward packets coming from specific sources to a different IP and the rest to the main one.
The routes can be automatically added using BGP/OSPF/etc and matched with a specific next-hop, once the route is added, a routing rule will be created by the project and will route the subnet under a different routing table. 

This could be used for a lot of other things as well - customize it as you will.

## Installation

1. Install Go on the build environment
2. Clone the project
3. Run `go get .` to download the dependencies
4. Update the [configuration file](https://github.com/amit177/go-syncrouterules/blob/main/config.toml) to match your environment
5. Run `make` to compile the program - change the ARCH using appropriate [environment variables](https://pkg.go.dev/cmd/go#hdr-Environment_variables) if needed
6. Upload the binary to the target machine
7. Run as root / with appropriate network capabilities


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.