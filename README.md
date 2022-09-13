# WaterMarker

Add a watermark to a whole bunch of photos.
For a set of photos I needed a variant with watermark for distribution, and a set without watermark for promotion. Instead of exporting the photos twice (once with watermark, once without) this tool allows a much faster method of creating a duplicate set of photos that is watermarked.

## How to use

In a single folder, include:

- a file called 'watermark.png'
- a folder called 'photos' containing jpg photos
- the addWatermark executable

```
$ ./addWatermark -h 	// To show options for command line arguments
$ ./addWatermark 	    // Create a folder 'watermarked' and populate with watermarked photos
```
