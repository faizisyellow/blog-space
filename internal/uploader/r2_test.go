package uploader

// import (
// 	"context"
// 	"net/http"
// 	"os"
// 	"testing"

// 	"github.com/charmbracelet/log"
// 	"github.com/joho/godotenv"
// )

// func TestR2Service(t *testing.T) {

// 	err := godotenv.Load("../../.env")
// 	if err != nil {
// 		log.Fatal("error loading .env file", "msg", err)
// 	}

// 	r2conf := struct {
// 		BucketName      string
// 		AccountId       string
// 		AccessKeyId     string
// 		AccessKeySecret string
// 	}{
// 		BucketName:      os.Getenv("R2_BUCKET_NAME"),
// 		AccountId:       os.Getenv("R2_ACCOUNT_ID"),
// 		AccessKeyId:     os.Getenv("R2_ACCESS_KEY"),
// 		AccessKeySecret: os.Getenv("R2_ACCESS_KEY_SECRET"),
// 	}

// 	r2client := NewR2Client(r2conf.BucketName, r2conf.AccountId, r2conf.AccessKeyId, r2conf.AccessKeySecret)

// 	t.Run("should successfully upload image from file system to R2 bucket", func(t *testing.T) {
// 		filePath := "./lizzy.jpeg"

// 		file, err := os.Open(filePath)
// 		if err != nil {
// 			t.Fatal("failed to open file:", err)
// 		}
// 		defer file.Close()

// 		sniff := make([]byte, 512)
// 		n, err := file.Read(sniff)
// 		if err != nil {
// 			t.Fatal("failed to read file header:", err)
// 		}
// 		contentType := http.DetectContentType(sniff[:n])

// 		_, err = file.Seek(0, 0)
// 		if err != nil {
// 			t.Fatal("failed to seek to start:", err)
// 		}

// 		filename := "lizzy.jpeg"

// 		err = r2client.UploadFile(context.TODO(), filename, file, contentType)
// 		if err != nil {
// 			t.Errorf("upload failed: %v", err)
// 		}
// 	})

// }
