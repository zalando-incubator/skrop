package filters

import (
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"io"
	"io/ioutil"
	"math"
)

type ImageFilter interface {
	CreateOptions(image *bimg.Image) (*bimg.Options, error)
}

func handleResponse(ctx filters.FilterContext, f ImageFilter) {
	rsp := ctx.Response()

	rsp.Header.Del("Content-Length")

	in := rsp.Body
	r, w := io.Pipe()
	rsp.Body = r

	go handleImageTransform(w, in, f)
}

func handleImageTransform(out *io.PipeWriter, in io.ReadCloser, f ImageFilter) error {
	defer func() {
		in.Close()
	}()

	imageByes, err := ioutil.ReadAll(in)

	if err != nil {
		return err
	}

	responseImage := bimg.NewImage(imageByes)

	options, err := f.CreateOptions(responseImage)

	if err != nil {
		return err
	}

	return transformImage(out, responseImage, options)
}

func transformImage(out *io.PipeWriter, image *bimg.Image, opts *bimg.Options) error {
	var err error

	defer func() {
		if err == nil {
			err = io.EOF
		}
		out.CloseWithError(err)
	}()

	newImage, err := image.Process(*opts)

	if err != nil {
		return err
	}

	_, err = out.Write(newImage)

	return err
}

func parseEskipIntArg(arg interface{}) (int, error) {
	if number, ok := arg.(float64); ok && math.Trunc(number) == number {
		return int(number), nil
	} else {
		return 0, filters.ErrInvalidFilterParameters
	}
}