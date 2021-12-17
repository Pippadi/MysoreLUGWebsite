# Dual-Booting Linux with OSX on a Core Duo MacBook 1,1
## By Prithvi Vishak
## Created July 13th, 2021
## Updated July 22nd, 2021
---

Recently, my parents made a trip to my grandparents' house, and they brought back my grandpa's old MacBook. It is a MacBook Pro 1,1 complete with an Intel Core Duo processor, 2GB of RAM, and an 80GB hard disk.
OSX Snow Leopard ran a little slow, and was obviously completely useless. It was time for the Linux treatment. Installing Linux on a modern-ish mac is fairly straightforward, so I thought this would be the same way.

*_Boy, was I wrong._*

Modern PCs use the Universal Extensible Firmware Interface (UEFI), the modern successor to BIOS. Intel-based macs do too, but as always, Apple thought different (did their own nonsense), and made their own implementation of it.
Relatively modern 64-bit Intel macs (without the T2 security chip) with their custom EFI boot Linux, even the live disks, just fine with pretty much no tweaking needed.
As it turns out, Apple's firmware on the MacBook Pro 1,1 is significantly different. It just won't boot Linux from a thumb drive.

This article from [AstroFloyd](https://astrofloyd.wordpress.com/2014/01/14/linux-only-installation-on-2006-macbook-using-refind/) worked great for a lot of people, but unfortunately not for me.
I couldn't get Linux to install without OSX, but I wanted to keep OSX around anyway, in case I'd ever need it. And dual-booting is always cooler than single-booting.
Anyway, this is a chronicle of my findings.


First off, I installed the [rEFInd boot manager](https://www.rodsbooks.com/refind/) by running the following (OSX doesn't seem to have `wget`).

```
prithvi $ curl -S webwerks.dl.sourceforge.net/project/refind/0.13.2/refind-bin-0.13.2.zip &gt; refind.zip
...
prithvi $ unzip refind.zip
...
prithvi $ cd refind-bin-0.13.2
prithvi $ sudo ./refind-install
```

When I rebooted, I was greeted by rEFInd. In Disk Utility, I shrank OSX's partition.

!assets/diskutil.png
!width="60%" height="auto" alt="OSX Disk Utility"
OSX is greedy.

After this, I flashed Zorin OS, my 32-bit distro of choice, to a thumb drive, and booted to the live environment.
I chose to create my partitions in GParted, because I like its interface.
I created only an EFI partition and root partition for Linux as shown, of course not touching OSX.

!assets/gparted.png
!width="80%" height="auto" alt="Disk partitioning"
Do you notice something different?

Next, I went through the installer's advanced partitioning, and set it to use the EFI and root partition we created, and nothing more.
I made sure it wasn't using OSX's EFI partition (which it was by default), and only its own.

Oddly enough, it did not use that EFI partition for anything. It seemed to boot in BIOS mode with the mac's Compatibility Support Module (CSM), which emulates BIOS.
Apparently, this is used for installing Windows with Bootcamp, which explains why I saw an entry called "Windows" in the boot menu (not rEFInd), which boots Zorin.

That was about it to make it boot. Of course, this was the result of multiple failed attempts, where I nuked OSX.
Reinstalling OSX was actually what took me the most time. After I installed rEFInd the first time, I tried installing openSUSE, which messed with OSX's EFI partition, leaving the system unbootable.
We had some Snow Leopard DVDs lying around, but the optical drive on the old mac was long dead. Thankfully, we had another mac with a working optical drive, from which I could use Disk Utility to flash the contents of the DVD to a thumb drive.
Booting and installing from the thumb drive worked fine. If you need it and live close by, I'd be happy to lend it to you :)

Most things worked out of the box for me on Zorin OS Lite 15.3 (Ubuntu 18.04 LTS), including display and keyboard brightness controls, keyboard and trackpad, WiFi, bluetooth, speakers, and microphone.
I can't say much about the battery, because mine has been long-dead.
The first thing that didn't work on Linux was the iSight webcam. I followed [this tutorial](https://unlockforus.com/linux-mint-17-2-rafaela-on-macbook-41-late-2007-isight-webcam/), which worked great for me.

The last thing to not work was suspend. When I tried suspending, I would see the system trying to suspend, but it would wake up immediately as if suspend failed.
This was indeed the case. Looking at `dmesg`, I found some errors saying that it wasn't able to report suspend state to the TPM module of all things.
Apparently, these 2006 Intel macs were the last macs to have TPMs. That sentence there was pretty much everything I could find about them online.
Either way, I had no need for TPM, so I tried blacklisting the driver. Finding that suspend did indeed work after doing so, I made the changes permanent by creating the file `/etc/modprobe.d/tpm-blacklist.conf`, with contents:

```
blacklist tpm
blacklist tpm_infineon
```

My problems didn't end there, unfortunately. On my third day of trying to get Linux to work, I tried starting the mac as usual. I heard the optical drive start and hard disk spin up.
Immediately after, I heard a click from the hard disk, and saw the activity/sleep light turn off. I eventually got it to boot by pressing-and-holding the power button till I saw the light flicker and heard the unit beep.
This happened for a few days, and disappeared as mysteriously as it came. According to this [Macrumors thread](https://forums.macrumors.com/threads/macbook-only-turns-on-when-power-button-is-held-down.1285597/), it's due to liquid damage on the top case.
Whatever it is, it's gone now, and it still works.

Anyway, as I was saying, Zorin OS only boots with BIOS. This is because the mac has a 32-bit EFI, and Ubuntu ISOs apparently don't come with 32-bit EFI support.
So, the installer could only install a BIOS-bootable system.
After looking online, I found that one could create an EFI-bootable Ubuntu ISO by copying some of the EFI files from a 32-bit Debian ISO (which does support 32-bit EFI boot out of the box).
Unfortunately, none of the ISOs I created worked (fine, I only tried it once, but it took a long time and I didn't want to do it again). So instead, I decided to install Debian directly.
I mean, what's cooler than 2 OSs? 3 OSs, obviously.

I first went for a Debian stable (10) live ISO with LXQt and non-free drivers. I got GRUB from EFI. The live environment didn't seem to work (more on this later), but the installer did.
It seemed to install okay, but when I tried booting to it, all I got was a black screen. Attributing it to an ancient kernel (4.19, why on earth do people use Debian stable?), I downloaded a Debian testing (11, kernel 5.10) weekly build ISO, which also seemed to install correctly.
Still a black screen after GRUB.

After yet more research, I found that this was a known issue with older Radeon cards (which this has), and could perhaps be due to a [lack of nonfree firmware](https://superuser.com/questions/1289728/fresh-installed-debian-9-is-bootin-into-black-screen).
This involved editing the boot options in GRUB, by adding `radeon.modeset=0` to the end of the line starting with `linux`.
Installing the nonfree firmware didn't help, but the boot option got me to a terminal, where `startx` started an LXQt session.
The reason I wasn't booting to a GUI seemed to be an issue with the display manager, `SDDM`, not starting.
I wasn't in the mood to debug it, so I downloaded a Debian stable netinstaller, and installed that.

Repeating the same steps to disable modesetting, I still got only a terminal. When I tried `startx`, it wouldn't start, saying that it needed modesetting.

*_Forget it._*

!assets/refind-rabbit-hole.jpg
!width="75%" height="auto" alt="rEFInd Menu"
Rabbit hole? What rabbit hole?

I knew graphics worked on Zorin OS, so it occured to me that I could _add_ EFI boot support from Debian by installing GRUB from there. I had an empty EFI partition and everything.
Following the instructions [here](https://wiki.debian.org/GrubEFIReinstall) didn't result in any errors.

Selecting Ubuntu with `grubia32.efi` from rEFInd got me GRUB. When I tried booting as usual, I got the same blank screen. The keyboard backlight came on, suggesting that the OS seemed to have started X, but wasn't displaying properly.
So, I added the `radeon.modeset=0` to Zorin OS's boot options in GRUB, and I _finally_ got a GUI!

Unfortunately, this came with issues of its own. Graphics performance was _abysmal_. I suppose that everything is somehow CPU-rendered, or the graphics card is not clocked as fast as it could be. Also, I can't change the display brightness anymore.
This is the stage the computer is at now. I suppose I will go back to booting Zorin with BIOS until the graphics card quirks are ironed out.

I continue to experiment and troubleshoot with this machine (I say "this" machine because I've used it to write most of this article). My long-term goal is to get it to boot Linux from a thumb drive directly, without rEFInd or the tool used by AstroFloyd.
That's all for now. I hope this article was vaguely useful, even though there's a dearth of better-researched, less rambling articles out there.

I'd like to give a shoutout to [Roderick Smith](https://www.rodsbooks.com), whose work is very detailed and informative. Not to mention that he created rEFInd, if you didn't notice.
His articles helped me understand a lot about the boot process and mac-related quirks.
