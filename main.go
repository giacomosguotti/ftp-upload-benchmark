package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	ftp "github.com/moov-io/go-ftp"
)

// Arguments:
//  -host: ftp server to call
//	-user: ftp server user
//  -pwd: ftp server pwd
//  -file: path of the file to upload
//  -n: number of file uploads to perform
//  -out: file path of output measurements
//  -dst: destination path on the FTP server

// Output csv file format:
// 1st column: start timestamp of the measurement
// 2nd column: elapsed time of the connection to the FTP server
// 3rd column: elapsed time of the file upload
// 4th column: elapsed time of the file upload + connection to the FTP server
// 5th column: error message (if any occurred during the upload process)

func main() {

	// *** Command-Line Flags ***
	hostPtr := flag.String("host", "ftp.localhost:21", "address of the ftp server")
	userPtr := flag.String("user", "root", "user for ftp server authentication")
	pwdPtr := flag.String("pwd", "password", "password for ftp server authentication")
	uploadFilePtr := flag.String("uploadFile", "uploadFile", "path of the file to upload")
	nPtr := flag.Int("n", 100, "number of uploads to perform")
	outPtr := flag.String("out", "results.csv", "output file path with measurements obtained")
	serverDstPtr := flag.String("dst", "/tmp/x", "destination path on the FTP server")

	flag.Parse()

	// Show command-line flags
	fmt.Println("")

	// Initialize the FTP client configuration
	clientConfig := ftp.ClientConfig{
		Hostname: *hostPtr,
		Username: *userPtr,
		Password: *pwdPtr,
	}

	// Open output CSV file for appending (create if not exists)
	outFile, err := os.OpenFile(*outPtr, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	defer outFile.Close()
	csvWriter := csv.NewWriter(outFile)
	defer csvWriter.Flush()

	// Write header if file is empty
	fileInfo, err := outFile.Stat()
	if err == nil && fileInfo.Size() == 0 {
		csvWriter.Write([]string{
			"start_timestamp",
			"open_elapsed_time_ms",
			"upload_elapsed_time_ms",
			"total_elapsed_time_ms",
			"error",
		})
		csvWriter.Flush()
	}

	// *** Start measurement campaign ***
	for i := 0; i < *nPtr; i++ {

		start := time.Now()
		startTimestamp := start.Format(time.RFC3339Nano)
		var openElapsedTime, uploadElapsedTime time.Duration
		var errMsg string

		// Open a connection to the FTP server using the provided configuration
		client, err := OpenFtpConnection(clientConfig)
		if err != nil {
			fmt.Println("Error connecting to FTP server:", err)
			openElapsedTime = time.Since(start)
			errMsg = err.Error()
			// Write metrics to CSV
			csvWriter.Write([]string{
				startTimestamp,
				fmt.Sprintf("%d", openElapsedTime.Milliseconds()),
				"", // upload time not measured
				fmt.Sprintf("%d", openElapsedTime.Milliseconds()),
				errMsg,
			})
			csvWriter.Flush()
			time.Sleep(1 * time.Second)
			continue
		}

		openElapsedTime = time.Since(start)

		err = UploadFtpFile(client, *uploadFilePtr, *serverDstPtr)
		if err != nil {
			fmt.Println("Error uploading file:", err)
			uploadElapsedTime = time.Since(start) - openElapsedTime
			errMsg = err.Error()
			client.Close()
			// Write metrics to CSV
			csvWriter.Write([]string{
				startTimestamp,
				fmt.Sprintf("%d", openElapsedTime.Milliseconds()),
				fmt.Sprintf("%d", uploadElapsedTime.Milliseconds()),
				fmt.Sprintf("%d", (openElapsedTime + uploadElapsedTime).Milliseconds()),
				errMsg,
			})
			csvWriter.Flush()
			time.Sleep(1 * time.Second)
			continue
		}

		uploadElapsedTime = time.Since(start) - openElapsedTime

		// Close the FTP client connection
		err = client.Close()
		if err != nil {
			fmt.Println("Error closing FTP client:", err)
			errMsg = err.Error()
			// Write metrics to CSV
			csvWriter.Write([]string{
				startTimestamp,
				fmt.Sprintf("%d", openElapsedTime.Milliseconds()),
				fmt.Sprintf("%d", uploadElapsedTime.Milliseconds()),
				fmt.Sprintf("%d", (openElapsedTime + uploadElapsedTime).Milliseconds()),
				errMsg,
			})
			csvWriter.Flush()
			time.Sleep(1 * time.Second)
			continue
		}

		// Write successful metrics to CSV
		csvWriter.Write([]string{
			startTimestamp,
			fmt.Sprintf("%d", openElapsedTime.Milliseconds()),
			fmt.Sprintf("%d", uploadElapsedTime.Milliseconds()),
			fmt.Sprintf("%d", (openElapsedTime + uploadElapsedTime).Milliseconds()),
			"",
		})
		csvWriter.Flush()

		time.Sleep(1 * time.Second)
	}

	fmt.Println("End")
}

// OpenFtpConnection opens a connection to the FTP server using the provided configuration.
// It returns the FTP client, and an error if any occurred.
// The function is responsible for creating a new FTP client and checking if the connection is reachable.
func OpenFtpConnection(clientCfg ftp.ClientConfig) (client ftp.Client, err error) {
	client, err = ftp.NewClient(clientCfg)

	if err != nil {
		return client, err
	}

	// Check if the FTP client is reachable
	if err := client.Ping(); err != nil {
		return client, err
	}
	return client, nil
}

// UploadFtpFile uploads a file to the FTP server using the provided client and file path. The file is uploaded to the specified destination path on the server.
// It returns an error if any occurred during the upload process.
// The function is responsible for opening the file to be uploaded and closing it after the upload is complete.
// It also handles the upload process and checks for any errors that may occur during the upload.
func UploadFtpFile(client ftp.Client, filePath string, destination string) (err error) {

	// Open the file to be uploaded
	fileData, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileData.Close()

	// Upload the file to the destination path
	err = client.UploadFile(destination, fileData)
	if err != nil {
		return err
	}

	return nil
}
