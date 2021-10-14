package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "will download a file from GCS",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := creds()
		if err != nil {
			return err
		}

		return RetrieveInputFile(context.Background(), filename, c, bucketName)
	},
}

var (
	filename   string
	bucketName string
)

func init() {
	downloadCmd.Flags().StringVarP(&filename, "file", "f", "", "filename on the GCS bucket")
	downloadCmd.Flags().StringVarP(&bucketName, "bucket", "b", "", "bucket name on the GCS")
	rootCmd.AddCommand(downloadCmd)
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

func RetrieveInputFile(ctx context.Context, filename string, creds []byte, bucketName string) error {
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
	objecthandler := bucket.Object(filename)
	reader, err := objecthandler.NewReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}
