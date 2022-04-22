# Sharing Your Screen as a Webcam
## Prithvi Vishak
## Created November 27th, 2021

---

Our school uses Microsoft Teams for online classes. Unfortunately, the IT admins blocked students from sharing their screens, to the dismay of both students and teachers.
I worked around the problem by redirecting screen contents to webcam video.

There seem to be many ways to do this, some easier than others. I had heard of [OBS](https://obsproject.com/), but didn't need a tool that big for regular old screen-sharing.
After some hunting around the internet, I pieced together a script that uses the [`v4l2loopback`](https://github.com/umlaeute/v4l2loopback) kernel module and [`ffmpeg`](https://ffmpeg.org/) to redirect the desired screen space to a webcam.

### My Script

That script can be found on [GitHub](https://github.com/Pippadi/ScreenToWebcam). To use it, install `v4l2loopback` and `ffmpeg`.
Both should be available from your distro's official repositories. Keep in mind that the package names may be slightly different, and to run a search if it says that packages by those names are not available.
Also install <code>git</code> if you want to clone the repository instead of downloading the [zip from GitHub](https://github.com/Pippadi/ScreenToWebcam/archive/refs/heads/main.zip).

Once you have the script, launch a terminal and navigate to the folder that contains the script. Run it with your desired settings, then hit `Ctrl-C` to stop.

```
./ScreenToWebcam  [options] [values]
E.g.:  ./ScreenToWebcam --mirror 1680x1050 1920x1080 <br>
-m,--mirror: Flip webcam feed horizontally<br>
-- Values --
InputSize:  Required field
            Size of the screen to be grabbed
            '1680x1050' in example given
OutputRes:  Optional field
            Resolution the input will be scaled to in the dummy webcam output
            Defaults to 1280x720
```

I run it like this, to capture 1920x1080 pixels:

```
./ScreenToWebcam.sh 1920x1080
```

Resolution is specified in the format `WIDTHxHEIGHT+WIDTH\_OFFSET,HEIGHT\_OFFSET`.
If you have two 1920x1080 monitors side-by-side, and you want to capture the right one, you would do:

```
./ScreenToWebcam.sh 1920x1080+1920,0
```

The resolution of the dummy webcam is set by 1280x720 by default, as this is the highest resolution that Teams and a couple other applications seem to accept.
The input screen size is scaled to the resolution of the dummy webcam, so if you are capturing an area with an aspect ratio different from 16:9, your webcam feed will appear stretched to fit the entire frame. You may want to set the output resolution to one with a matching aspect ratio.
Additionally, if the captured portion of your screen has a higher resolution than 1280x720, text may appear small and illegible. It is advisable to increase the size of whatever you are sharing, or even set your monitor to 1280x720 resolution for the duration of your sharing.

You will also be prompted for your password to run `modprobe` as root when starting and stopping the script, because the script instantiates the `v4l2loopback` kernel module at start and removes when stopped, _every time the script is run_. Also keep in mind that the script runs _only_ on Xorg/X11, and not on Wayland.

!assets/webcam-preview.png
!width="50%" height="auto" alt="ScreenToWebcam preview in Teams"
This is *without* the `--mirror` flag in use. While the preview is shown mirrored here in Teams, people in meetings tell me that they see it *correctly*, and flipped when I use `--mirror`.

### Doing it yourself

As mentioned earlier, you need the `v4l2loopback` module to create a dummy webcam device. Then, we use `ffmpeg` with the `x11grab` device to grab a portion of our screen.
You can instantiate `v4l2loopback` with `modprobe` every time before running ffmpeg like my script does, or you could [add an entry to /etc/modprobe.d](https://askubuntu.com/questions/1245212/how-do-i-automatically-run-modprobe-v4l2loopback-on-boot).

Here's the command you would run before the script, to create a dummy webcam with the name *ScreenToWebcam* at `/dev/video10` \(Keep in mind that this will remove any other instances of `v4l2loopback` you might have running, and assumes that you don't already have a webcam using `/dev/video10`\):

```
sudo /sbin/modprobe v4l2loopback devices=1 exclusive_caps=1 card_label=ScreenToWebcam video_nr=10
```

And the `ffmpeg` command to grab a resolution of 1920x1080, at an offset of 1920x0 pixels from the top-left corner of the screen, with no mirroring:

```
ffmpeg -f x11grab -video_size 1920x1080 -i $DISPLAY+1920+0 -vf scale=1280x720,format=yuv420p -r 15 -c:a copy -f v4l2 /dev/video10
```

When you're done, hit `Ctrl-C` to terminate `ffmpeg`, and run

```
sudo /sbin/modprobe -r v4l2loopback
```

This is all the information I'll leave you with, since this is pretty much all I know. Feel free to refer to the [script directly](https://github.com/Pippadi/ScreenToWebcam/blob/main/ScreenToWebcam.sh). Happy hacking!
