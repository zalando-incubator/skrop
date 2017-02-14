package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"gopkg.in/h2non/bimg.v1"
)

func TestNewCropByHeight(t *testing.T) {
	name := NewCropByHeight().Name()
	assert.Equal(t, "cropByHeight", name)
}

func TestCropByHeight_Name(t *testing.T) {
	c := cropByHeight{}
	assert.Equal(t, "cropByHeight", c.Name())
}

func TestCropByHeight_CreateOptions(t *testing.T) {
	c := cropByHeight{height: 400, cropType: North}
	image := imagefiltertest.LandscapeImage()
	options, _ := c.CreateOptions(image)

	assert.Equal(t, 1000, options.Width)
	assert.Equal(t, 400, options.Height)
	assert.Equal(t, true, options.Crop)
	assert.Equal(t, bimg.GravityNorth, options.Gravity)
}

func TestCropByHeight_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewCropByHeight, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  true,
	}, {
		Msg:  "one arg",
		Args: []interface{}{400.0},
		Err:  false,
	}, {
		Msg:  "two args",
		Args: []interface{}{400.0, North},
		Err:  false,
	}, {
		Msg:  "more than 2 args",
		Args: []interface{}{400.0, 200.0, North},
		Err:  true,
	}})
}
