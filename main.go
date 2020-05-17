package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/egymgmbh/go-prefix-writer/prefixer"
	"github.com/gdrive-org/gdrive/auth"
	"github.com/gdrive-org/gdrive/drive"
	"github.com/patrickdappollonio/env"
)

var (
	quietMode    = env.GetBoolean("FBACKUP_QUIET")
	acctFile     = env.GetDefault("FBACKUP_ACCOUNT_FILE", "")
	targetFolder = env.GetDefault("FBACKUP_FOLDER", "")
	destFolder   = env.GetDefault("FBACKUP_DESTINATION", "")
	everyN       = env.GetDefault("FBACKUP_EVERY", "30m")
)

func main() {
	logger := log.New(os.Stdout, "[fbackup] ", log.LstdFlags)
	if quietMode {
		logger.SetOutput(ioutil.Discard)
	}

	acctFile, err := checkFile(acctFile, "account file", false)
	if err != nil {
		logger.Fatalf("Unable to initialize fbackup: %s", err.Error())
	}

	logger.Printf("Setting account file to: %s", acctFile)

	targetFolder, err := checkFile(targetFolder, "target folder", true)
	if err != nil {
		logger.Fatalf("Unable to initialize fbackup: %s", err.Error())
	}

	logger.Printf("Setting target folder to: %s", targetFolder)

	if everyN == "" {
		logger.Fatal("Unable to initialize fbackup: duration is empty")
	}

	if destFolder == "" {
		logger.Fatal("Unable to initialize fbackup: destination folder ID is empty")
	}

	every, err := time.ParseDuration(everyN)
	if err != nil {
		logger.Fatalf("Unable to parse fbackup duration %q: %s", everyN, err.Error())
	}

	logger.Printf("Setting ticker to: %s", every.String())

	client, err := auth.NewServiceAccountClient(acctFile)
	if err != nil {
		logger.Fatalf("Unable to log in to Google Drive with service account provided: %s", err.Error())
	}

	logger.Printf("Configured service account with account file: %s", acctFile)

	gdrive, err := drive.New(client)
	if err != nil {
		logger.Fatalf("Unable to instantiate a Google Drive client using account provided: %s", err.Error())
	}

	logger.Printf("Google Drive client configured, retrieving folder information for ID: %s", destFolder)

	if err := gdrive.Info(drive.FileInfoArgs{Id: destFolder, Out: ioutil.Discard}); err != nil {
		logger.Fatalf("Unable to get information about the destination folder in Google Drive: make sure it exists and you can write to it. Error: %s", err.Error())
	}

	logger.Printf("Folder information gathered, started ticking...")

	ticker := time.NewTicker(every)
	defer ticker.Stop()

	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, os.Kill)

	logger.Println("Initiating first zip and backup process")

	if err := backup(logger, gdrive, targetFolder, destFolder); err != nil {
		logger.Fatalf("Unable to zip and backup: %s", err.Error())
	}

	logger.Printf("First backup succeeded. Waiting for next tick in ~%s.", every)

	for {
		select {
		case <-sigquit:
			logger.Println("Received quit signal, closing application.")
			return

		case <-ticker.C:
			if err := backup(logger, gdrive, targetFolder, destFolder); err != nil {
				logger.Fatalf("Unable to zip and backup: %s", err.Error())
			}

			logger.Printf("Backup succeeded. Waiting for next tick in ~%s.", every)
		}
	}
}

func backup(logger *log.Logger, gd *drive.Drive, targetFolder, destFolder string) error {
	logger.Println("Uploading folder contents:", targetFolder, " -- Drive folder ID:", destFolder)

	wr := prefixer.New(logger.Writer(), func() string {
		return "[fbackup] " + time.Now().Format("2006/01/02 15:04:05") + " [gdrive] "
	})

	uploadArgs := drive.UploadSyncArgs{
		Out:              wr,
		Progress:         ioutil.Discard,
		Path:             targetFolder,
		RootId:           destFolder,
		DeleteExtraneous: true,
		Timeout:          20 * time.Second,
		Resolution:       drive.KeepLocal,
		Comparer:         &dateComparer{},
	}

	if err := gd.UploadSync(uploadArgs); err != nil {
		return fmt.Errorf("unable to upload folder %q: %w", targetFolder, err)
	}

	return nil
}
