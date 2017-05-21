/*
<!--
Copyright (c) 2017 Christoph Berger. Some rights reserved.

Use of the text in this file is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

Use of the code in this file is governed by a BSD 3-clause license that can be found
in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = ""
description = "Control a Digispark in Go via Gobot and LittleWire"
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2017-05-20"
draft = "true"
domains = ["Robotics"]
tags = ["digispark", "microcontroller", "littlewire", "gobot"]
categories = ["Project"]
+++

### Summary goes here

The Digispark is perhaps as small (feature-wise) as a microcontroller board can get. This is a short writedown about my first experiences with controlling this board through Go code, using Gobot and LittleWire.

<!--more-->

## A tiny microcontroller board

It was after the last Munich Gophers Meetup when a few of us went to a local bar to discuss about Go, life, and hardware. Yes, hardware. From foldable USB keyboards to Raspberry Pi Zero W to some incredibly small microcontroller board called Digispark.

This one caught my attention. The Digispark is a board built around an ATtiny85 microcontroller chip, and features:

* Incredible six (!) I/O ports
* A whopping 8 KB of flash memory (that's over eight thousand bytes, folks!)
* Blazingly fast 16.5MHz system clock

This tiny board connects directly to an USB port of a host computer and can be equipped with Arduino scripts (with restrictions of course), but this is not an Arduino blog, so here comes the Go part.

## Controlling a Digispark board from Go code

With a few more (software) ingredients, the I/O ports can be controlled from Go code running on a PC/Mac/LinuxBox/etc:

* LittleWire
* The `micronucleus` cli app
* Gobot

[LittleWire](http://littlewire.cc) is basically a script that adds USB communication capabilities to the Digispark. With an accompanying library at the USB host's end, the board's I/O ports can be remote-controlled from an app running on the host.

[`micronucleus`](https://github.com/micronucleus/micronucleus) is a bootloader that comes preinstalled with Digisparks. (It takes about 2KB from the microcontroller's flash memory, leaving about 6KB memory available for scripts.) It has a companion CLI app that is needed for installing LittleWire on the device.

[Gobot](https://gobot.io/) is a Go robotics framework that connects Go apps to a large array of electronic devices, from the little Digispark to Arduino, Raspberry Pi, and even Quadrocopters.

Now let's put these pieces together.


## The steps

Here are the steps I took to make the onboard status LED blink. (A blinking LED is the "Hello World" of microcontrollers.)


### Step 1: Install the micronucleus CLI app

I compiled the `micronucleus` CLI tool from its source in two steps:


#### Clone the source

A simple `git clone` fetches the source code from GitHub:

    git clone https://github.com/micronucleus/micronucleus


#### Compile the command

After this, I cd'ed into `micronucleus/commandline` and ran `make`:

```
commandline $  [master|✔] make
Building library: micronucleus_lib...
gcc -I/usr/local/Cellar/libusb-compat/0.1.5_1/include -L/usr/local/Cellar/libusb-compat/0.1.5_1/lib -lusb -Ilibrary -O -g -D MAC_OS -c library/micronucleus_lib.c
```

On my Mac, I got a couple of warnings but as the `micronucleus` command turned out to work fine, I consider these warnings as benign. On my Lubuntu netbook, `make` finished without any warnings.


### Step 2: Install LittleWire

Again, only two steps are required.

#### Download littlewire_13.hex

The "firmware" script LittleWire is available at [the LittleWire homepage](http://littlewire.cc/downloads.html)

Note: While writing this, I noticed that the site seems not consistently available, but luckily you can access an [archived version](https://web.archive.org/web/20161024004122/http://littlewire.cc/downloads.html) at the Web Archive.


#### Install LittleWire on the device

Then I ran `micronucleus`:

    commandline $  [master|✔] ./micronucleus --run ../../LittleWire/littlewire_v13.hex

(I had LittleWire in a directory next to the Git workspace of micronucleus; if you want to repeat these steps, adapt the path as needed.)

`micronucleus` then asked to plug in the device into an USB port. (The Digispark, when powered on, listens on the USB port for some seconds, and if the host does not send anything, the device disables the USB port, because two of the I/O pins have a double use as USB channels, and disabling the USB port makes these two ports available for hardware wired to the board. This is why the board must be plugged in only when the host app is ready.)

I plugged in the board, and `micronucleus` uploaded the LittleWire firmware in a few seconds.


### Step 3: Install Gobot

For this, I simply followed the instructions from the [Gobot documentation](https://gobot.io/documentation/platforms/digispark/).

For Mac and Linux, some USB libraries need to be installed. On the Mac, this is a one-liner if you have [Homebrew](https://brew.sh) installed.

    brew install libusb && brew install libusb-compat

On Ubuntu (and surely also on other Debian-based distributions), a simple `apt-get` does the trick.

    sudo apt-get install libusb-dev

Now I was ready to fetch Gobot and install the Digispark platform package. (If you are new to Go: the ellipsis at the end of the `go get` line advise the `go get` tool to also download all subprojects, even if the top project does not import them. The `-d` is also important as this prevents the `go get` tool from installing anything at this point. The only thing we want to install is the Digispark package.)

      go get -d -u gobot.io/x/gobot/...

The next step installs the `digispark` library. I first ran into an error ("undefined: lw") until I realized that I had `cgo` disabled. On a Bash shell, the prefix `CGO_ENABLED=1` ensures that go install uses `cgo`.

      CGO_ENABLED=1 go install gobot.io/x/gobot/platforms/digispark

While this worked on the Mac, on Linux I got a "package libusb was not found" error message. I double-checked that `libusb-dev` was installed on the system. Um-hm...


## The code
*/

// ## Imports and globals
package main

/*
## How to get and run the code

Step 1: `go get` the code. Note the `-d` flag that prevents auto-installing
the binary into `$GOPATH/bin`.

    go get -d github.com/appliedgo/TODO:

Step 2: `cd` to the source code directory.

    cd $GOPATH/src/github.com/appliedgo/TODO:

Step 3. Run the binary.

    go run TODO:.go


## Odds and ends
## Some remarks
## Tips
## Links


**Happy coding!**

*/
