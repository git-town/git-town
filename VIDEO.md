# Creating animations of the development work

First, install [gource](https://gource.io). It seems to work best on Windows.

Preview the generated movie:

```fish
gource --load-config .gource.conf
```

Create a video file:

```fish
gource --load-config .gource.conf -o - | ffmpeg -y -r 60 -f image2pipe -vcodec ppm -i - -vcodec libx264 -preset veryslow -pix_fmt yuv420p -threads 0 -bf 0 -crf 18 git-town.mp4
```
