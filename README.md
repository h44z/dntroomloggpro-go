# DNT RoomLogg PRO

A simple program that can read the temperature and humidity values from the DNT RoomLogg PRO base station using USB.

Project inspired by https://juergen.rocks/blog/articles/elv-raumklimastation-rs500-raspberry-pi-linux.html (https://github.com/juergen-rocks/raumklima).

## Setup
Compile the go code using the provided Makefile:
```shell
make
```

## Installation and Usage
In order to execute the compiled program without administrator rights, use the provided udev rule:
```shell
sudo cp scripts/99-hid.rules /etc/udev/rules.d/99-hid.rules
```

Now you can use the `reader` binary to read the stats from  the DNT RoomLogg PRO base station.