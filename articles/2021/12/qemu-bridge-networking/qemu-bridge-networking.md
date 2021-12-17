# Bridge Networking with Qemu and Libvirt
## By Prithvi Vishak
## Created December 17th, 2021
---

### Backstory

My dad's work requires him to use Windows. The catch here, though, is that he's often dealing with multiple conflicting versions of software, and in general cluttering up a single system with it all (_dependency hell_ for you package maintainers out there).
Despite (because of?) having used Windows for so long, he took to Linux quite nicely, and worked in Windows VMs in VirtualBox (he did this on Windows at one point too).
He also tried out Libvirt with Qemu, and liked the speed it offered relative to VirtualBox.

My dad likes his VMs on our local network. VirtualBox made bridging to your local network a trivial task just a switch or two to flip, but that didn't seem to be the case with Libvirt and virt-manager.
Eventually, he ended up switching back to Windows for another project, but he's expressed interest in switching back to Linux, this time using Libvirt with bridge networking.
He assigned me the task of getting it to work, since I already used Libvirt for my VM needs.

Network bridging seems to only cooperate with physical wired interfaces. My dad and I use WiFi all the time because my mom won't let us run a cable.
In my surfings of the interwebs, I discovered that some networking whizzes work around this limitation by creating dummy ethernet interface thingamabobs and having those hooked up to their physical wireless interface.

Yes, I too had a hard time following that, but then came across routing, which seems to allow for bridge-like behavior over a wireless connection.
Long (and I mean _really_ long) story short, that didn't work for me either (probably because I don't understand it yet, and the tutorials are all older than Gangam Style).
If I make progress on this front, I'll update this article.

Recently, I chanced upon [Adam Williamson's wonderful article](https://www.happyassassin.net/posts/2014/07/23/bridged-networking-for-libvirt-with-networkmanager-2014-fedora-21/) on getting bridging working quickly and easily from a GUI (because I'm a normie).
This article describes the _*exact*_ same process, with some screenshots for mostly mine and my dad's reference, along with a few minor changes adjusting for modern versions of apps.

### What I did

I decided to work around the wired-only limitation of bridges by digging out an old router, configuring it as an extender, placing it at my desk, and plugging my computer into it via ethernet (through a USB dongle - thanks a *lot*, Dell).
Next, I disabled the automatically created ethernet profile in Settings under the Networks tab.

!assets/disable-network.png
!width="70%" height="auto" alt="Disable the network"
For the grandpas out there

Then, I installed the NetworkManager Connections Editor package (called `NetworkManager-connection-editor` on openSUSE). Weirdly, the ability to create bridges seems to have been _removed_ from the GNOME settings app, as per other articles I went through.
I can see why they get so much heat from users for arbitrarily removing features.

Anyway, I launched the app (naturally called "Advanced Network Configuration" instead of "Connection Editor"), clicked *+*, chose to create a *Bridge*, then hit *Create*.
I use DHCP for IP address allocation on my home network. If you don't, you may want to edit things in the IPV4 section.
Other than that, no changes needed to be made in the current window (I mean, if you use IPV6, you probably don't need this article).

Under the *Bridge* pane, I hit *Add* in the *Bridged Connections* section, and chose Ethernet in the dropdown. In the *Ethernet* pane, select your physical ethernet device in the drop-down.

!assets/left-to-right.png
!width="90%" height="auto" alt="Steps"
Left to right in chronological order

Hitting *Save*, *Create*, and whatever other buttons I needed to make things final, I didn't see a new connection show up in the Settings app's Networks section, as described in above article.
For whatever reason, there doesn't seem to be a way to activate a connection in the "Advanced" Network Connections tool, so I went to the terminal, and ran this:

```
~> nmcli con show
NAME                     UUID                                  TYPE       DEVICE  
wlp6s0                   6a2352b6-b29f-4c8a-a45c-b235235accf0  wifi       wlp6s0  
Bridge connection 1      03504422-7d40-45b8-9f4f-564a23452354  bridge     bridge0 
bridge0 port 1           08332b0b-c7f2-4638-9449-17e98235423c  ethernet   --      
...
~> nmcli con up 08332b0b-c7f2-4638-9449-17e98235423c
```

The first command was for getting the UUID of the bridge _port_ (not the bridge itself). The second was to activate the connection. This is the host's connection, and doesn't have anything to do with VMs yet.
Unlike what the article describes, I still don't see an entry in the GNOME Settings app.

Well, that's about it.
In my VM's configuration, I set the NIC to the new bridge device, and I was on my way!

!assets/vm-settings.png
!width="65%" height="auto" alt="VM Settings"

Now all I need to do is convince my mom to let me run a cable from the router to get rid of the extender... the hardest part of this whole endeavor.
