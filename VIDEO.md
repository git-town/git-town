# Creating animations of the development work

First, install [gource](https://gource.io). It seems to be the most solid on Windows.

Preview the generated movie:

```bash
gource --load-config .gource.conf
```

Create a video file:

```bash
gource --load-config .gource.conf | ffmpeg -y -r 60 -f image2pipe -vcodec ppm -i - -vcodec libx264 -preset veryslow -pix_fmt yuv420p -threads 0 -b:v 600k git-town.mp4
```
