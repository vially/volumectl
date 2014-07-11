# volumectl - CLI utility to control the volume


## Features

![notification bubble](https://wiki.ubuntu.com/Sound?action=AttachFile&do=get&target=notification.jpg)

`volumectl` displays a nice notification bubble when [notify-osd](https://launchpad.net/ubuntu/+source/notify-osd) is installed.


## Usage Examples

##### Increase volume (by default the increment is 2%):
`$ volumectl up` 

##### Increase volume with 7%:
`$ volumectl set 7%+`

##### Decrease volume (by default the decrement is 2%):
`$ volumectl down`

##### Toggle mute:
`$ volumectl toggle`

##### Set volume to 50%:
`$ volumectl set 50%`


## Installation

### Arch Linux
`volumectl` is available from the [AUR](https://aur.archlinux.org/packages/volumectl/)

### Building from source

#### Requirements
 * [Go](http://golang.org/doc/install) is required to build this package from source
 * `alsa-utils` is needed because the `amixer` utility is called by `volumectl`
 * [notify-osd](https://launchpad.net/ubuntu/+source/notify-osd) is required for displaying the notifications

#### Download and install package
`go get github.com/vially/volumectl`


## Available Commands
    USAGE:
       volumectl command [arguments]
    
    COMMANDS:
       up          increase volume (with 2%)
       down        decrease volume (with 2%)
       mute        mute volume
       unmute      unmute volume
       toggle      toggle mute
       set         set volume to a specific value
       help, h     Shows a list of commands or help for one command
       
    GLOBAL OPTIONS:
       --version, -v    print the version
       --help, -h       show help
