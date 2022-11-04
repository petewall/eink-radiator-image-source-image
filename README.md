# eInk Radiator Image Source: Image

![CI](https://ci.petewall.net/api/v1/teams/main/pipelines/eink-radiator/jobs/test-image-source-image/badge)

Generates an image from an existing image.

```bash
image --config config.json --height 300 --width 400
```

## Configuration

The configuration is the image source, which must be a publically accessible URL, the scaling algorithm (see below) and the background color (if required).

```json
{
    "src": "http://example.com/surfing.jpg",
    "scale": "cover",
    "background": {
        "color": "purple"
    }
}
```

Possible options for `scale`:

* `resize` - Resize the image to fit the desired resolution. May lead to distortions.
* `contain` - Resize the image so the whole image fits inside the new resolution. May show some background, which will use the `background.color` configuration.
* `cover` - Resize the image so the smallest dimension fits inside the new resolution frame. May crop out some of the original image.
