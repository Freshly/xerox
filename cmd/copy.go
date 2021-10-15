package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copies (xerox is a copy machine, get it?) a file from GCS bucket to your local machine.",
	Long: `Copies (xerox is a copy machine, get it?) a file from GCS bucket to your local machine.
Set XEROX_CREDS_BASE64 env to configure the credentials for your project.
By default, your current gcloud profile will be used for authentication.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := creds()
		if err != nil {
			return err
		}

		return copyFileFromBucket(context.Background(), filename, bucketName, c)
	},
}

var (
	filename   string
	bucketName string
)

func init() {
	copyCmd.Flags().StringVarP(&filename, "file", "f", "", "filename on the GCS bucket")
	copyCmd.Flags().StringVarP(&bucketName, "bucket", "b", "", "bucket name on the GCS")
	rootCmd.AddCommand(copyCmd)
}

func creds() ([]byte, error) {
	encryptedCreds := os.Getenv("XEROX_CREDS_BASE64")
	if encryptedCreds == "" {
		return nil, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(encryptedCreds)
	if err != nil {
		return nil, fmt.Errorf("decoding creds: %v", err)
	}
	return decoded, nil
}

func copyFileFromBucket(ctx context.Context, filename string, bucketName string, creds []byte) error {
	var client *storage.Client
	var err error

	if creds == nil {
		client, err = storage.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("error creating GCS client: %v", err)
		}
	} else {
		client, err = storage.NewClient(ctx, option.WithCredentialsJSON(creds))
		if err != nil {
			return fmt.Errorf("error creating GCS client from the creds: %v", err)
		}
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	objectHandler := bucket.Object(filename)
	reader, err := objectHandler.NewReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	attributes, err := objectHandler.Attrs(ctx)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	fmt.Println(attributes.Size/1000, "KB to download")
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " Downloading..."
	s.FinalMSG = "\nDone!\n"
	s.Start()
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}
	s.Stop()

	return nil
}
