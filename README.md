# route53-ddns

A simple CLI written in Go that will update Route53 A Records with your public DNS.

This uses Viper for config and accepts an array of DNS records to update.
Check out the route53-ddns.yaml file in the deb directory for the format.
This uses the IPify API to get the public IP address.
This project uses zerolog with lumberjack for logging and log rotation.
You can change the logging level by setting the env variable LOGLEVEL.

## TODO

* Add AWS SDK verbose logging to get request-ids
* Add ability to set logging level in viper
* Get the CircleCI build working
* Create a deb package
