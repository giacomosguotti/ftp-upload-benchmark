# ftp-upload-benchmark

This Go program measures the time taken to connect to an FTP server and upload a file, logging the results to a CSV file. It is designed for repeated, automated measurements of FTP upload performance.

## Features

- Connects to an FTP server using configurable credentials.
- Uploads a specified file to a specified path on the FTP server.
- Measures and logs:
  - Start timestamp
  - Connection time
  - Upload time
  - Total time (connection + upload)
  - Error messages (if any)
- Appends results to a CSV file after each iteration.
- Command-line configurable.

## Requirements

- Go 1.18 or newer
- [github.com/moov-io/go-ftp](https://github.com/moov-io/go-ftp) library

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/giacomosguotti/ftp-upload-benchmark
cd ftp-upload-benchmark
go mod tidy
go build -o ftp-upload-benchmark main.go
```

## Usage

```bash
./ftp-upload-benchmark \
  -host <ftp_host:port> \
  -user <ftp_user> \
  -pwd <ftp_password> \
  -uploadFile <local_file_path> \
  -n <number_of_iterations> \
  -out <output_csv_file> \
  -dst <remote_ftp_path>
```

### Arguments

- `-host` (string): FTP server address and port (e.g., `ftp.localhost:21`)
- `-user` (string): FTP username
- `-pwd` (string): FTP password
- `-uploadFile` (string): Path to the local file to upload
- `-n` (int): Number of upload iterations to perform (default: 100)
- `-out` (string): Output CSV file path for measurements (default: `results.csv`)
- `-dst` (string): Destination path on the FTP server (e.g., `/tmp/testfile.txt`)

### Example

```bash
./ftp-upload-benchmark \
  -host ftp.localhost:21 \
  -user root \
  -pwd password \
  -uploadFile ./testfile.txt \
  -n 100 \
  -out results.csv \
  -dst /tmp/testfile.txt
```

## Output

The program appends results to the specified CSV file after each iteration. The CSV columns are:

1. `start_timestamp`: Start time of the measurement (RFC3339Nano)
2. `open_elapsed_time_ms`: Time to connect to the FTP server (milliseconds)
3. `upload_elapsed_time_ms`: Time to upload the file (milliseconds)
4. `total_elapsed_time_ms`: Total time (connection + upload) in milliseconds
5. `error`: Error message (if any occurred during the upload process)

Example CSV output:

```
start_timestamp,open_elapsed_time_ms,upload_elapsed_time_ms,total_elapsed_time_ms,error
2024-06-01T12:00:00.123456789Z,50,120,170,
2024-06-01T12:00:02.123456789Z,48,118,166,connection refused
...
```

## Notes

- The program waits 1 second between each iteration.
- If an error occurs during connection, upload, or closing, it is logged in the CSV file.
- The output CSV file is created if it does not exist, and a header is written if the file is empty.

## Troubleshooting

- Ensure the FTP server is reachable from your machine.
- Make sure the FTP credentials and paths are correct.
- The program requires permission to read the local file and write to the output CSV.
- If you see errors in the CSV, check the error column for details.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---
