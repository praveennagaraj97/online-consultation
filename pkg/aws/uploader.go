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
	filePath string,
	fileName string,
	formFileKey string,
	width, height uint64,
	ch chan *S3UploadChannelResponse,
) {

	file, err := ctx.FormFile(formFileKey)
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

	originalPath := fmt.Sprintf("%s/original/%s.%s", filePath, fileName, fileExtension)
	blurPath := strings.Replace(originalPath, "original", "blur", 1)

	_, err = a.UploadAsset(bytes.NewBuffer(buffer), originalPath, &fileType)
	if err != nil {
		ch <- &S3UploadChannelResponse{
			Result: nil,
			Err:    err,
		}
		return
	}

	blurBuffer, err := utils.CreateBlurDataForImages(buffer, 1, int(width)/2, int(height)/2)
	if err != nil {
		ch <- &S3UploadChannelResponse{
			Result: nil,
			Err:    err,
		}
		return
	}

	_, err = a.UploadAsset(bytes.NewBuffer(blurBuffer), blurPath, &fileType)
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
			OriginalSrc:       fmt.Sprintf("%s/%s", a.options.S3_PUBLIC_ACCESS_BASEURL, originalPath),
			BlurDataURL:       fmt.Sprintf("%s/%s", a.options.S3_PUBLIC_ACCESS_BASEURL, blurPath),
			Width:             width,
			Height:            height,
		},
		Err: err,
	}
}
