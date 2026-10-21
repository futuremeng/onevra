package gallery

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// Load and decode once, including failures, rather than retrying every frame.
// MyGo resolves bundle resources and source/resources for direct go run.
var galleryImage = sync.OnceValues(func() (*ui.Bitmap, error) {
	data, err := galleryImageData()
	if err != nil {
		return nil, err
	}
	return ui.DecodeBitmap(data)
})

// galleryImageData is the encoded avatar, which notifications send as is.
var galleryImageData = sync.OnceValues(func() ([]byte, error) {
	dir, err := mygo.App.Path(mygo.PathResources)
	if err != nil {
		return nil, err
	}
	return readGalleryImage(dir)
})

func readGalleryImage(dir string) ([]byte, error) {
	return os.ReadFile(filepath.Join(dir, "images", "avatar.jpeg"))
}

func loadGalleryImage(dir string) (*ui.Bitmap, error) {
	data, err := readGalleryImage(dir)
	if err != nil {
		return nil, err
	}
	return ui.DecodeBitmap(data)
}

func controlImage(c *ui.Context, _ *controlState, variant string) {
	bitmap, err := galleryImage()
	if err != nil {
		components.Text(c, tr(c, "unableToLoadAvatarJpeg")).TextColor(c.Theme().TextMuted)
		return
	}
	image := components.Image(c, bitmap).Size(180, 120).Label(msg(c).AuthorAvatarS(tr(c, variant)))
	switch variant {
	case "Cover":
		image.Fit(ui.Cover)
	case "Stretch":
		image.Fit(ui.FillBox)
	case "Rounded":
		image.Size(120, 120).Radius(24)
	case "Grayscale":
		image.Grayscale()
	}
}
