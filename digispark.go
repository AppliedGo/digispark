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
title = "Controlling a Digispark board"
description = "Control a Digispark microcontroller in Go via Gobot and LittleWire"
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2017-05-21"
draft = "false"
domains = ["Robotics"]
tags = ["digispark", "microcontroller", "littlewire", "gobot", "hardware", "iot"]
categories = ["Project"]
+++

The Digispark is perhaps as small as a microcontroller board for DIY electronics can get. This is a short writedown about my first experiences with controlling this board through Go code, using Gobot and LittleWire.

<!--more-->

## A tiny microcontroller board

It was after the latest Munich Gophers Meetup when a few of us went to a local bar to talk about Go, life, and hardware. Yes, hardware. From foldable USB keyboards to Raspberry Pi Zero W to some incredibly small microcontroller board called Digispark.

![A Digispark clone](digispark.jpg)

This one caught my attention. (Ok, the board is nothing new, but I did not hear of it before.) The Digispark is a board built around an ATtiny85 microcontroller chip and features -

* Incredible **six (!) I/O ports**
* A whopping **8 KB of flash memory** (that's over eight thousand bytes, folks!)
* A blazingly fast **16.5MHz system clock**

This tiny board connects directly to an USB port of a host computer and can be equipped with Arduino scripts (with restrictions of course, it is a much smaller controller than the one on the Arduino boards), but this is not an Arduino blog, so here comes the Go part.


## Controlling a Digispark board from Go code

With a few more (software) ingredients, the I/O ports can be controlled from Go code running on a PC/Mac/LinuxBox/etc. These are:

* LittleWire
* The `micronucleus` cli app
* Gobot

[LittleWire](http://littlewire.cc) is basically a script that adds USB communication capabilities to the Digispark. With an accompanying library at the USB host's end, the board's I/O ports can be remote-controlled from an app running on the host.

[`micronucleus`](https://github.com/micronucleus/micronucleus) is a bootloader that comes preinstalled with Digisparks. (It needs a little space in the microcontroller's flash memory, leaving about 6KB memory available for scripts.) It has a companion CLI app that is needed for installing LittleWire on the device.

[Gobot](https://gobot.io/) is a Go robotics framework that connects Go apps to a large array of electronic devices, from the little Digispark to Arduino, Raspberry Pi, and even Quadrocopters.

![A first success](digisparkblink.jpg)
*(Image: A first success)*

And here is how I used all this to make the onboard status LED blink, and a servo motor move.

## The steps


### Step 1: Install the micronucleus CLI app

I compiled the `micronucleus` CLI tool from its source in two steps:


#### Clone the source

A simple `git clone` fetches the source code from GitHub:

    git clone https://github.com/micronucleus/micronucleus


#### Compile the command

After this, I cd'ed into `micronucleus/commandline` and ran `make`:

```
commandline $ make
Building library: micronucleus_lib...
gcc -I/usr/local/Cellar/libusb-compat/0.1.5_1/include -L/usr/local/Cellar/libusb-compat/0.1.5_1/lib -lusb -Ilibrary -O -g -D MAC_OS -c library/micronucleus_lib.c
```

On my Mac, I got a couple of warnings but as the `micronucleus` command turned out to work fine, I consider these warnings as benign. On my Lubuntu netbook, `make` finished without any warnings.


### Step 2: Install LittleWire

Again, only two steps are required.

#### Download littlewire_13.hex

The "firmware" script LittleWire is available at [the LittleWire download page](http://littlewire.cc/downloads.html)

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

While this worked on the Mac, on Linux I got a "package libusb was not found" error message. Apparently the error was raised by a tool named `pkg-config`. I wasn't able to figure out the reason for this, so my Linux-based test ended here.


## The code

After everything was in place, I wrote the below code to make the onboard LED blink and let a servo move from 0 to 180 degrees in steps of 45 degrees. This video shows the outcome:

<iframe src="https://player.vimeo.com/video/218360443" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>


But now to the code:

*/

//
package main

import (
	"log"
	"time"

	// Besides the `gobot` package, we also need the `gpio` package for controlling
	// the pins, and the `digispark` package for talking to the Digispark board via
	// USB.
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/digispark"
)

func main() {
	// An adaptor provides the connection to a device.
	digispark := digispark.NewAdaptor()

	// The `gpio` package contains drivers for various sensor and actor hardware.
	// Here, we use the LED and servo drivers.
	// Each driver needs the adaptor to use and the number of the pin to control.
	led := gpio.NewLedDriver(digispark, "1")
	servo := gpio.NewServoDriver(digispark, "5")

	// `work` is the main routine that a Gobot robot executes. We will pass this func
	// to `NewRobot`.
	work := func() {
		// `Every` plans actions in fixed intervals. It does not wait for previous
		// actions to complete before triggering the next action.
		ledTicker := gobot.Every(2*time.Second, func() {
			for i := 0; i < 10; i++ {
				// `After` schedules an action after a specific time. It does not wait.
				gobot.After(time.Duration(i)*100*time.Millisecond, func() {
					// Here we switch the LED on and off.
					led.Toggle()
				})
			}
		})
		gobot.After(6*time.Second, func() {
			// `Every` returns a `*Ticker` that we stop after a while.
			ledTicker.Stop()
		})

		servoTicker := gobot.Every(7*time.Second, func() {
			// `Move` moves the servo into the given position
			// (between 0 and 180 degrees).
			servo.Move(45)
			gobot.After(1*time.Second, func() {
				servo.Center()
			})
			gobot.After(2*time.Second, func() {
				servo.Move(135)
			})
			gobot.After(3*time.Second, func() {
				// 180 is the theoretical maximum, but this makes my servo try moving too far. 170 works just fine.
				servo.Move(170)
			})
			gobot.After(4*time.Second, func() {
				// My servo thinks "0" is -10 degrees, so I had to adjust the "0" position a bit.
				servo.Move(10)
			})
		})
		// Let's stop this after 30 seconds.
		gobot.After(30*time.Second, func() {
			servoTicker.Stop()
		})

	}

	// Now we can create a "robot" using the items we defined above.
	robot := gobot.NewRobot("AppliedGoBot",
		[]gobot.Connection{digispark},
		[]gobot.Device{led, servo},
		work,
	)

	// Finally we just need to start our robot.
	err := robot.Start()
	if err != nil {
		log.Println(err)
	}
}

/*
## How to get and run the code

This code is only useful if you have a Digispark at hand, but in any case, here are the installation instructions.

Step 1: `go get` the code. Note the `-d` flag that prevents auto-installing
the binary into `$GOPATH/bin`.

    go get -d github.com/appliedgo/digispark

Step 2: `cd` to the source code directory.

    cd $GOPATH/src/github.com/appliedgo/digispark

Step 3. Run the binary.

    go run digispark.go

(Stop with `Ctrl-C`.)


## Conclusion

Controlling hardware can be quite fun! I already envision a couple of neat (and useless) things I could create:

* An analog voltage meter that shows the number of emails in my inbox
* An array of LED's for monitoring the status of remote servers
* A big, red panic button that triggers any kind of emergency action (instant server shutdown, for example)
* and more!

A big part of the fun surely is the Gobot framework. It certainly deserves a closer look! Maybe I write another post when/if I have enough time and some more hardware to test out.


### Caveats

* Ultra-low-cost devices like the Digispark often have no means of electrical protection. So if you connect the wrong things to the I/O ports (or the right things in a wrong way), **you might not only fry the microcontroller board but maybe also the USB chip on your computer's mainboard, and maybe even the mainboard itself.** So be careful about what you are doing, and always check your wiring twice! Also, a USB hub can lower the risk of damaging your mainboard if you fry your microcontroller board.
* The Digispark is cheap, and some clones are even cheaper. (The Digispark is open hardware.) Mine is a clone, as it has a "rev3" tag although no such revision exists for the original Digispark. I got it for EUR 1,95, and at this price tag, don't expect high quality. For example, my board exposes a strange behavior when the servo is moved for the first time: The on-board LED then stops working until I unplug the device. I am not sure if other boards expose this problem as well. Maybe this particular board is a lemon but maybe the problem is in the hardware specs, or maybe the clones I have are just of low quality. (I assume the latter.)


### Can it get cheaper than that?

If you want to control devices, e.g. an LC display, without any additional controller hardware, hop over to [this post](https://dave.cheney.net/2014/08/03/tinyterm-a-silly-terminal-emulator-written-in-go) on Dave Cheney's blog and learn how to connect an LC display to an unused monitor port (HDMI, DVI, VGA - either of these works).


**Happy coding!**

*/
