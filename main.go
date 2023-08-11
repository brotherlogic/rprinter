package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	pbp "github.com/brotherlogic/printqueue/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func localPrint(ctx context.Context, lines []string) error {
	//os.Create("home/simon/print.txt")
	handle, err := os.CreateTemp("", "printdetails")
	if err != nil {
		return fmt.Errorf("unable to create temporary file: %w", err)
	}
	for _, line := range lines {
		handle.WriteString(fmt.Sprintf("%v\n", line))
	}
	err = handle.Sync()
	if err != nil {
		return fmt.Errorf("unable to sync file: %w", err)
	}
	err = handle.Close()
	if err != nil {
		return fmt.Errorf("unable to close file: %w", err)
	}

	cmd := exec.Command("lp", fmt.Sprintf("%v", handle.Name()))
	output := ""
	out, err := cmd.StdoutPipe()

	if err != nil {
		return fmt.Errorf("error in the now resolved actual stdout: %w", err)
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
		return fmt.Errorf("error in running command %w (output: %v)", err, output)
	}

	log.Printf("Printed %v", handle.Name())

	return err
}

func runReceiptPrint() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := grpc.Dial("print.brotherlogic-backend.com:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		err = localPrint(ctx, job.GetLines())
		if err != nil {
			return err
		}

		_, err = client.Ack(ctx, &pbp.AckRequest{
			Id:      job.GetId(),
			AckType: pbp.Destination_DESTINATION_RECEIPT,
		})
		if err != nil {
			return fmt.Errorf("ACK error on %v -> %w", job.GetId(), err)
		}
	}

	return nil
}

func main() {
	err := runReceiptPrint()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
