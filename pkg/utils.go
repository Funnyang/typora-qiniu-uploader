package pkg

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"os"

	_ "image/gif"
	_ "image/png"

	"github.com/nfnt/resize"
)

// PathExists check if a path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// CompressImg 压缩图片
func CompressImg(r io.Reader) (io.Reader, int64, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		if err == image.ErrFormat {
			return r, 0, nil
		}
		return nil, 0, err
	}
	level := []uint{3240, 2160, 1920, 1280, 1024}
	buf := bytes.NewBuffer(nil)
	for i := 0; i < len(level); i++ {
		m := resize.Resize(level[i], 0, img, resize.Lanczos3)
		// jpg
		if err = jpeg.Encode(buf, m, &jpeg.Options{Quality: 90}); err != nil {
			return nil, 0, err
		}
		// 2M
		if buf.Len() < 2097152 {
			break
		} else {
			buf = bytes.NewBuffer(nil)
		}
	}
	return buf, int64(buf.Len()), nil
}
