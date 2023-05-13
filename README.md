# WaterMarker

Add a watermark to a whole bunch of photos.
For a set of photos I needed a variant with watermark for distribution, and a set without watermark for promotion. Instead of exporting the photos twice (once with watermark, once without) this tool allows a much faster method of creating a duplicate set of photos that is watermarked.

# Development
Install ![Wails](https://wails.io/docs/gettingstarted/installation/), dependencies include:
- Golang
- NPM

Change the code in app.go and the frontend folder (uses Svelte) as you see fit.

Then run `./build.sh` if you're on linux, on windows I just execute the command in the `build.sh` file:
- `wails build -platform windows/amd64 -o WaterMarker_Windows.exe`