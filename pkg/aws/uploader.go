package awspkg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	"github.com/praveennagaraj97/online-consultation/utils"
)

type S3UploadChannelResponse struct {
	Result *interfaces.ImageType
	Err    error
}

func (a *AWSConfiguration) UploadImageToS3(ctx *gin.Context,
	fileKey string,
	ch chan *S3UploadChannelResponse,
	fileName string,
	filePath string) {

	file, err := ctx.FormFile(fileKey)
	if err != nil {
		ch <- &S3UploadChannelResponse{
			Result: nil,
			Err:    err,
		}
		return
	}

	fileType := file.Header.Get("Content-Type")

	if !strings.Contains(fileType, "image") {
		ch <- &S3UploadChannelResponse{
			Result: nil,
			Err:    errors.New("Provided file is not acceptable"),
		}
		return
	}

	// Read the file buffer.
	multiPartFile, err := file.Open()
	if err != nil {
		ch <- &S3UploadChannelResponse{
			Result: nil,
			Err:    err,
		}
		return
	}
	buffer, err := io.ReadAll(multiPartFile)
	if err != nil {
		ch <- &S3UploadChannelResponse{
			Result: nil,
			Err:    err,
		}
		return
	}

	defer multiPartFile.Close()

	fileExtensionClips := strings.Split(file.Filename, ".")
	fileExtension := fileExtensionClips[len(fileExtensionClips)-1]
	fileKeyName := fmt.Sprintf("%s.%s", fileName, fileExtension)
	originalPath := filePath + "/original"
	blurPath := filePath + "/blur"

	_, err = a.UploadAsset(bytes.NewBuffer(buffer), originalPath, fileKeyName, &fileType)
	if err != nil {
		ch <- &S3UploadChannelResponse{
			Result: nil,
			Err:    err,
		}
		return
	}

	blurBuffer, err := utils.CreateBlurDataForImages(buffer, 1, 10, 10)
	if err != nil {
		ch <- &S3UploadChannelResponse{
			Result: nil,
			Err:    err,
		}
		return
	}

	_, err = a.UploadAsset(bytes.NewBuffer(blurBuffer), blurPath, fileKeyName, &fileType)
	if err != nil {
		ch <- &S3UploadChannelResponse{
			Result: nil,
			Err:    err,
		}
		return
	}

	ch <- &S3UploadChannelResponse{
		Result: &interfaces.ImageType{
			OriginalImagePath: originalPath,
			BlurImagePath:     blurPath,
		},
		Err: err,
	}
}
