package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	pbp "github.com/brotherlogic/printqueue/proto"
	"google.golang.org/grpc"
)

func localPrint(ctx context.Context, lines []string) error {
	//os.Create("home/simon/print.txt")
	handle, err := os.CreateTemp("rprint", "printdetails")
	if err != nil {
		return fmt.Errorf("Unable to create temporary file: %w", err)
	}
	for _, line := range lines {
		handle.WriteString(fmt.Sprintf("%v\n", line))
	}
	err = handle.Sync()
	if err != nil {
		return fmt.Errorf("Unable to sync file: %w", err)
	}
	err = handle.Close()
	if err != nil {
		return fmt.Errorf("Unable to close file: %w", err)
	}

	cmd := exec.Command("lp", handle.Name())
	output := ""
	out, err := cmd.StdoutPipe()

	if err != nil {
		return fmt.Errorf("Error in the now resolved actual stdout: %w", err)
	}

	if out != nil {
		scanner := bufio.NewScanner(out)
		go func() {
			for scanner != nil && scanner.Scan() {
				output += scanner.Text()
			}
			out.Close()
		}()
	}

	cmd.Start()
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("Error in running command %w (output: %v)", err, output)
	}

	return err
}

func runReceiptPrint() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := grpc.Dial("printer.brotherlogic-backend.com:8000")
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	client := pbp.NewPrintServiceClient(conn)
	jobs, err := client.RegisterPrinter(ctx, &pbp.RegisterPrinterRequest{
		Id:           hostname,
		ReceiverType: pbp.Destination_DESTINATION_RECEIPT,
	})

	if err != nil {
		return err
	}

	for _, job := range jobs.GetJobs() {
		localPrint(ctx, job.GetLines())
	}

	return nil
}
