<img src="https://raw.githubusercontent.com/kolja/skrop/e86f88b7d1452c7b07781e5152aea05f2c5c342c/skropodile.svg?sanitize=true"
     width="350"
     height="350"
     alt="skropodile">

# Skrop &nbsp; [![Build Status](https://travis-ci.org/zalando-stups/skrop.svg?branch=master)](https://travis-ci.org/zalando-stups/skrop) [![codecov](https://codecov.io/gh/zalando-stups/skrop/branch/master/graph/badge.svg)](https://codecov.io/gh/zalando-stups/skrop) [![Go Report Card](https://goreportcard.com/badge/github.com/zalando-stups/skrop)](https://goreportcard.com/report/github.com/zalando-stups/skrop) [![MicroBadger Size](https://img.shields.io/microbadger/image-size/skrop/skrop.svg?style=flat)](https://hub.docker.com/r/skrop/skrop/) [![Docker Pulls](https://img.shields.io/docker/pulls/skrop/skrop.svg?style=skrop)](https://hub.docker.com/r/skrop/skrop/) [![Docker](https://img.shields.io/badge/docker-skrop%2Fskrop-blue.svg?style=flat)](https://hub.docker.com/r/skrop/skrop/) [![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/zalando-stups/skrop/master/LICENSE)
Skrop is a media service based on [Skipper](https://github.com/zalando/skipper) and the [vips](https://github.com/jcupitt/libvips) library.

## Usage

In order to be able to use Skrop, you have to be familiar with how
[Skipper](https://github.com/zalando/skipper) works.

* [Install](./docs/INSTALL.md)
   * [Build from sources](./docs/INSTALL.md#build-from-sources)
     * [Install dependencies](./docs/INSTALL.md#install-dependencies)
   * [Using Docker](./docs/INSTALL.md#using-docker)
   * [Using Heroku](./docs/INSTALL.md#using-heroku)

### Getting started

### Run Skrop
```
go run cmd/skrop/main.go -routes-file eskip/sample.eskip -verbose
```

### Test

```
make all
```

To test if everything is configured correctly you should open in your browser
```
http://localhost:9090/images/big-ben.jpg
```
and the resized version
```
http://localhost:9090/images/S/big-ben.jpg
```

Here is how a route from the `sample.eskip` file would look like:

```
small: Path("/images/S/:image")
  -> modPath("^/images/S", "/images")
  -> longerEdgeResize(800)
  -> "http://localhost:9090";
```

What it means, is that when somebody does a `GET` to `http://skrop.url/images/S/myimage.jpg`,
Skrop will call `http://localhost:9090/images/myimage.jpg` to retrieve
the image from there, resize it, so that its longer edge is 800px and return
the resized image back in the response.


## Filters
Skrop provides a set of filters, which you can use within the routes:

* **longerEdgeResize(size)** — resizes the image to have the longer edge as specified, while at the same time preserving the aspect ratio
* **crop(width, height, type)** — crops the image to have the specified width and height the type can be "north", "south", "east" and "west"
* **cropByHeight(height, type)** — crops the image to have the specified height
* **cropByWidth(width, type)** — crops the image to have the specified width
* **resize(width, height, opt-keep-aspect-ratio)** — resizes an image. Third parameter is optional: "ignoreAspectRatio" to ignore the aspect ratio, anything else to keep it
* **addBackground(R, G, B)** — adds the background to a PNG image with transparency
* **convertImageType(type)** — converts between different formats (for the list of supported types see [here](https://github.com/h2non/bimg/blob/master/type.go)
* **sharpen(radius, X1, Y2, Y3, M1, M2)** — sharpens the image (for info about the meaning of the parameters and the suggested values see [here](http://www.vips.ecs.soton.ac.uk/supported/current/doc/html/libvips/libvips-convolution.html#vips-sharpen))
* **width(size, opt-enlarge)** — resizes the image to the specified width keeping the ratio. If the second arg is specified and it is equals to "DO_NOT_ENLARGE", the image will not be enlarged
* **height(size, opt-enlarge)** — resizes the image to the specified height keeping the ratio. If the second arg is specified and it is equals to "DO_NOT_ENLARGE", the image will not be enlarged
* **blur(sigma, min_ampl)** — blurs the image (for info see [here](http://www.vips.ecs.soton.ac.uk/supported/current/doc/html/libvips/libvips-convolution.html#vips-gaussblur))
* **imageOverlay(filename, opacity, gravity, opt-top-margin, opt-right-margin, opt-bottom-margin, opt-left-margin)** — puts an image onverlay over the required image
* **transformByQueryParams()** - transforms the image based on the request query parameters (supports only crop for now) e.g: localhost:9090/images/S/big-ben.jpg?crop=120,300,500,300.
* **cropByFocalPoint(targetX, targetY, aspectRatio, minWidth)** — crops the image based on a focal point on both the source as well as on the target and desired aspect ratio of the target. TargetX and TargetY are the definition of the target image focal point defined as relative values for both width and height, i.e. if the focal point of the target image should be right in the center it would be 0.5 and 0.5. This filter expects two PathParams named **focalPointX** and **focalPointY** which are absolute X and Y coordinates of the focal point in the source image. The fourth parameter is optional; when given the filter will ensure that the resulting image has at least the specified minimum width if not it will crop the biggest possible part based on the focal point.

### About filters
The eskip file defines a list of configuration. Every configuration is composed by a route and a list of filters to
apply for that specific route.
The finalizeResponse() needs to be added at the end, because it triggers the last transformation of the image.

Because of performance, most of the filters don't not trigger a transformation of the image, but if possible it is merged with
the result of the previous filters. The image is actually transformed every time the filter cannot be merged with the
previous one e.g. both edit the same attribute and also at the end of the filter chain by the finalizeResponse filter.

## Metadata
By default metadata are kept in the processed images. If you are not interested in metadata and 
you want them stripped from all the images that are processed, you can add the following 
environment variable to the running system:

```
STRIP_METADATA=TRUE
``` 

## Continuous Integration

We are using [_Travis CI_](https://travis-ci.org/zalando-stups/skrop) for this project.

### GitHub Token

In order for _Travis CI_ to perform its duty, a valid _GitHub_ token must be encrypted in
the [Travis configuration file](./.travis.yml).

If this token needs to change, here is how to set it up:
- Log on to _GitHub_ on the account you want _Travis CI_ to impersonate when performing operation on the repository
- Go to _GitHub_, in your [Personal access tokens](https://github.com/settings/tokens) configuration
- Click on _Generate new token_
- Check the `public_repo` scope
- Click _Generate token_

We will now encrypt that token with the [_travis CLI_](https://docs.travis-ci.com/user/encryption-keys/). Make sure it
is installed.

```bash
travis encrypt GITHUB_AUTH=the_token_copied_from_github
```

This will encrypt it in a way that **only** _Travis CI_ can decrypt. Your personal token is then quite safe.

The output should look like this.

```
  secure: "someBASE64value"
```

Take the output (the complete YAML key as it appears) and replace, in the [`.travis.yaml` file](./.travis.yml), this
line with the output of the previous command.

```yaml
env:
- secure: "someBASE64value"
```

## Versioning

This project uses [semantic versioning](https://semver.org/).

The patch-version (3rd digit) is bumped up automatically at every merge to master (by [_Travis CI_](/.travis.yml)).

### Increment the patch version
This is done automatically by _Travis CI_. Nothing special to do here. Example: merging when latest tag is `v3.23.291`
will automaticall tag a version `v3.23.292`.

### Increment the minor version
Since _Travis CI_ only automatically increases the patch-version, we need to manually pre-tag with the new version we
want.

Scenario:
- actual version is `v3.23.291`
- tag one of the commit on your branch with the new version you want.
  - it is **important** that the patch version be `-1`, since it will be incremented automatically by _Travis CI_.
  - `git tag v3.24.-1 && git push --tags`
- open the pull request.
- after merge, _Travis CI_ will tag automatically the right final version `v3.24.0`.
- delete the temporary manual tag
  - `git tag -d v3.24.-1 && git push --tags`

### Increment the major version
Since _Travis CI_ only automatically increases the patch-version, we need to manually pre-tag with the new version we
want.

Scenario:
- actual version is `v3.23.291`
- tag one of the commit on your branch with the new version you want.
  - it is **important** that the patch version be `-1`, since it will be incremented automatically by _Travis CI_.
  - `git tag v4.0.-1 && git push --tags`
- open the pull request.
- after merge, _Travis CI_ will tag automatically the right final version `v4.0.0`.
- delete the temporary manual tag
  - `git tag -d v4.0.-1 && git push --tags`


## License

[MIT](LICENSE)
