package textures

import (
	"errors"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"

	e "github.com/jsbento/PhxEngine/core/errors"
	"github.com/jsbento/PhxEngine/utils"
)

type TextureFormat string

const (
	Texture2D TextureFormat = "TEXTURE_2D"
)

type Texture struct {
	handle uint32
	width  int32
	height int32
	source string
	format TextureFormat
}

func NewTexture(source string, format TextureFormat) (*Texture, error) {
	texture := &Texture{
		source: source,
		format: format,
	}
	switch format {
	case Texture2D:
		if err := texture.loadTexture2D(); err != nil {
			return nil, err
		}
	}

	return texture, nil
}

func (t *Texture) loadTexture2D() error {
	img, info, err := utils.LoadImage(t.source)
	if err != nil {
		return err
	}
	t.width, t.height = int32(info.Width), int32(info.Height)
	internalFormat := gl.RGBA8
	dataFormat := gl.RGBA
	if info.Channels == 3 {
		internalFormat = gl.RGB8
		dataFormat = gl.RGB
	}
	var data []byte
	switch info.Format {
	case utils.ImageJPEG:
		d, err := utils.EncodeJPEG(img)
		if err != nil {
			return err
		}
		data = d
	case utils.ImagePNG:
		d, err := utils.EncodePNG(img)
		if err != nil {
			return err
		}
		data = d
	default:
		return errors.New("invalid image type")
	}

	var texture uint32
	gl.CreateTextures(gl.TEXTURE_2D, 1, &texture)
	if err := e.CheckGLError(); err != nil {
		return err
	}
	gl.TextureStorage2D(texture, 1, uint32(internalFormat), t.width, t.height)
	if err := e.CheckGLError(); err != nil {
		return err
	}

	gl.TextureParameteri(texture, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TextureParameteri(texture, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TextureParameteri(texture, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TextureParameteri(texture, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TextureSubImage2D(
		texture,
		0,
		0,
		0,
		t.width,
		t.height,
		uint32(dataFormat),
		gl.UNSIGNED_BYTE,
		unsafe.Pointer(unsafe.SliceData(data)),
	)
	if err := e.CheckGLError(); err != nil {
		return err
	}

	return nil
}

func (t *Texture) Bind(slot uint32) {
	gl.BindTextureUnit(slot, t.handle)
}

func (t *Texture) Cleanup() {
	gl.DeleteTextures(1, &t.handle)
	t.handle = 0
}
