# A Linux Setup for School Computers
## By Eeshaan K. B. and Prithvi Vishak
## Created February 25th, 2022
---

Our school's computers run Windows 7. Their Ivy Bridge processors don't support Windows 10, and just aren't great to use. It's time to put Linux on them.
The aim of this article is to propose how we would configure a Linux system as a school computer, specifically as a near drop-in replacement for a Windows client.
It also aims to address some concerns that surround such a project, from the IT department's perspective.

### The school's configuration

Our school computers are all on a single LAN. Each one has the same users (one for each grade) and applications installed.

Each user on these clients gets access to a file share hosted on a server accessible to all clients. For example, the `Grade 1` user gets access to a partticular share that contains all the files that the first graders need access too.
These shares generally contain folders for each student to save their data. This ensures that a student is able to sit at any station in the computer lab, and still have access to their work.

There certainly are finer details, but these are the salient points I wanted to address here.

### My plans to achieve this

First, I needed a file share server that is equivalent to the one at school. Since I don't know exactly how it's set up, I'm assuming it just hosts a bunch of SMB shares, each with different privileges.

Next, I needed to set up the Linux clients themselves. My plan was to start with one configured exactly the way I want it, and then clone it to other machines.

### Configuring the file server

Since I had a Fedora VM lying idle, I thought I'd use that as my SMB server.
The [official documentation](https://docs.fedoraproject.org/en-US/quick-docs/samba/) for setting up SMB shares is very useful, and suggested reading before proceeding with this section.

I first added system users for each SMB share I intend to add (`grade1` and `grade2` for the sake of this exercise), then passwords for SMB shares from the users.

```
sudo useradd grade1
sudo useradd grade2
sudo smbpasswd -a grade1
sudo smbpasswd -a grade2
```

For simplicity, I'm using their home directories (`/home/gradeX`) as the shares. Now to set some SELinux stuff:

```
sudo semanage fcontext --add --type "samba_share_t" /home/grade1
sudo restorecon -R /home/grade1
sudo semanage fcontext --add --type "samba_share_t" /home/grade2
sudo restorecon -R /home/grade2
```

Finally, I added the shares to the `samba` configuration.

```
[grade1]
        comment = Grade 1 Files
        path = /home/grade1
        writeable = yes
        browseable = yes
        public = yes
        create mask = 0644
        directory mask = 0755
        write list = user

[grade2]
        comment = Grade 2 Files
        path = /home/grade2
        writeable = yes
        browseable = yes
        public = yes
        create mask = 0644
        directory mask = 0755
        write list = user
```

And that was it. I could access the shares from my host machine without a hitch, when using the username `grade1` or `grade2` and the associated SMB password. Painless.

Communication between my VMs and host was possible due to my VMs being bridged onto my home network. Read [this article]("../../../2021/12/qemu-bridge-networking") on how to do that.

### Configuring the Linux Client

I chose Linux Mint for this, since it's user-friendly, has long-term support, and backed by a large group of people (both the Mint and Ubuntu communities).
I also chose Mint specifically because it had a very interesting-looking OEM install mode, where I could install the system as usual, and make whatever post-install configurations I want (like what OEMs do before shipping computers to customers) in a temporary account.

Technically, this can all be done without the OEM account, but I thought it would be cool to give the school a nice out-of-the-box experience with everything preinstalled.

First, I booted to the OEM installer and installed the system. From the OEM temporary account, I removed stuff a school computer wouldn't need, like Thunderbird and Timeshift, and installed some software like Python and Scratch. More on apps later.

```
sudo apt purge Thunderbird Timeshift ...
sudo apt autoremove
sudo apt update && sudo apt upgrade
sudo apt install idle-python3.9 scratch ...
```

Next, I added two users, `grade1`, and `grade2`, from the Users and Groups settings.
`grade1`'s UID and GID are both 1000, and `grade2`'s 1001.

Again, for simplicitiy's sake, I decided to mount the SMB shares as their home directories.
I added this to `/etc/fstab` for the shares at `192.168.142.205` (my Fedora VM) to be mounted with the users' permissions at their home directories.

```
//192.168.142.205/grade1 /home/grade1 cifs username=grade1,password=windows,uid=1000,gid=1000 0 0
//192.168.142.205/grade2 /home/grade2 cifs username=grade2,password=isterrible,uid=1001,gid=1001 0 0
```

I rebooted, removing the OEM temporary account. I set up an Admin account on first boot. That was all.
