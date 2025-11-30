# Netmonk Prime Device Monthly XLSX Reporting

Script to export monthly monitoring summary data to XLSX format.

## Project Overview
- **Configuration:**
  - Loaded by `config/config.go`
- **Entry Point:** `cmd/main.go`

## Main Features
CLI based applicaition to generate XLSX summary reporting by module and customer. Summary data are from bec-ga-network API.

## Running Services
### 1. Build the Binary
Run the following command to compile the Go application into a binary:
```
./scripts/build_linux.sh
```
For other operating system please check scripts folder
### 2. Configure the Service
Create a configuration file in the `./_bin/conf` directory. The application supports both `.yaml` and `.env` formats. If the configuration file is not specified, the program will search for the configuration in the OS environment variables.
### 3. Start the Service
Once the configuration is set and migrations have completed successfully, start the service using:
```
./_bin/prime-monthly-summary-reporting -config "./_bin/conf/cfg.yaml"
```