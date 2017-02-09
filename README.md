# Kernel Build Dashboard

This tool (i.e. `kbdashboard`) is used to configure and manage building process
of multiple linux kernels. It is written in Golang.

[![Build Status](https://travis-ci.org/choueric/kbdashboard.svg?branch=master)](https://travis-ci.org/choueric/kbdashboard)

# Overview

Developpers, especially in embedded system, may usually need to modify, build 
and test more than one linux kernels in various projects. Some may handle
different versions using different toolchains, some may one kernel version but 
different build directories based on different configuratios. In these cases, 
they have to remember various configurations and use differnt kinds of commands.

To simplify this matter, `kbdashboard`, acting as the wrapper,  handle all these
complicated details. It has some advantages as below:

- Run in any directory without changing to the direcotry of kernel sources.
- Use individual building directory without affecting the kernel source tree.
- Simple commands to perform various actions from starting up to installing.
- Configure easily by using the json format configuration file.
- Colorful shell output.

# Build this Tool

First you need a Golang build environment, which is easy to 
[install](https://golang.org/doc/install).

```
$ make
$ sudo make install
```

Then user will get the executable `kbdashboard` being installed into the 
direcotry `/usr/local/bin` and the bash completion script being installed
into `/etc/bash_completion.d`. The Makefile is very simple and can be modified
easily.

# How to Use this Tool

To start quickly:

```sh
$ kbdashboard edit
$ kbdashboard config def
$ kbdashboard config
$ kbdashboard build
$ kbdashboard install
```

More details are below.

## Create a Profile

The configuration file of `kbdashboard` is the core element. Use command:

```sh
$ kbdashboard edit
```
to create and modify the tool's configuration file. It is in JSON format.
If it is the first time to run this tool, the configuration file 
`$HOME/.config/kbdashboard/config.json` does not exist. The tool will create
a template one and open it using `vim`. The template is:

```json
{
	"editor": "vim",
	"current": 0,
	"profile": [
	{
		"name":"demo",
		"src_dir":"/home/user/kernel",
		"arch":"arm",
		"cross_compile":"arm-eabi-",
		"target":"uImage",
		"output_dir":"./_build",
		"defconfig":"at91rm9200_defconfig",
		"dtb":"at91rm9200ek.dtb",
		"mod_install_dir":"./_build/mod",
		"thread_num":4
	}
	]
}
```

In this template, there are some global options, such as `editor`, `current`, 
and only one profile with the name `demo`.

```
- editor  : Specify text editor which will be invoked when `edit` command is executed.
- current : Current profile index. If no speficy profile in `build` and `config` command,
            this index profile will be used.
``` 

The `profile` array can contain many profiles, each one is a build object. A 
profile has many options, which can be seen in the `demo` profile, but only two
are mandatory:

```
- name    : profile name.
- src_dir : directory path of kernel source.
```

Other options can be empty and will be the default value during building kernel.

```
- arch            : architecture, corresponding to `ARCH` of kernel build command.
- cross_compile   : cross compiler, corresponding to `CROSS_COMPILE` of kernel 
                    build command.
- target          : target of building kernel image.
- output_dir      : output build directory, corresponding to `O` of kernel build 
                    command. It is relative to the `src_dir` if not absolute.
- defconfig       : default configuration when start up.
- dtb             : the name of target DTB file.
- mod_install_dir : module install directory, corresponding to `INSTALL_MOD_PATH`
                    of kernel build command. It is relative to the `src_dir` if
		            not absolute.
- thread_num      : number of thread used to compile, corresponding to '-j' option.
```

After mofifying, use command `list` to see the result:

```sh
$ kbdashboard list
```

Use command `choose` to set it as the default profile:

```sh
$ kbdashboard choose demo
```

Use `list` command to see whether `demo` is marked by red asterisk.
From now on, it is the profile wich which other commands will deal.

## Configure Kernel

Before comiling the kernel, it has to do the default configuration:

```sh
$ kbdashboard config def
```

It will use `defconfig` in profile to generate `.config` file in the building
directory.

To do manually in menuconfig, use command:

```sh
$ kbdashboard config
or
$ kbdashboard config menu
```

To save configuration, use command:

```sh
$ kbdashboard config save
```

It will execute `savedefconfig`.

## Compile Kernel

To get the kernel image, use command:

```sh
$ kbdashboard build
or
$ kbdashboard build image
```

To compile driver modules and install to `mod_install_dir`, use command:

```sh
$ kbdashboard build modules
```

To get DTB file and install to `output_dir`, use command:

```sh
$ kbdashboard build dtb
```

## Install

Use `install` command to execute installation script which is writted by users:

```sh
$ kbdashboard install
```

At first, because there must be no installation script for this profile, the
above command will create an empty script and open it using the specified editor.

After user writting the new script, execute it by using the command again.

If user want to modify the script, use the command:

```sh
$ kbdashboard edit install
```

Using an user-defined install script is common in embedded development. But for
building kernel for PC, using `make` command with `install` and `modules_install`
targets is more useful.

# Details of Commands

This part may be out of date. See output of `help` for updated information.

## help

Use command `help` or empty command to see the help message:

```sh
$ kbdashboard help
or
$ kbdashboard
```

`help` can also show usage of specific command. For example:

```sh
$ kbdashboard help build
```

## list

List the profiles. The current profile is marked by the red asterisk symbol.

sub-command:

- `-v`: show all information of profiles.

## choose

Choose the current profile by name or index. for example:

```
$ kbdashboard choose 0
$ kbdashboard choose test

```

## edit

Edit the tool's configuration file with text-editor specified by `editor`.

sub-commands:

- `profile` : edit the tool's file of profiles. It is the default sub-command.
- `install`: edit the current profile's installation script.

## config

Do the configuration work of the current profile.

sub-command:

- `menu`: invoke `menuconfig`. It is the default sub-command.
- `def` : invoke `xxxx_defcofig` which specified by `defcofig` of the profile 
          configuration.
- `save`: invoke `savedefconfig`.

## build

Build various targets for kernel of current profile. 

- `image`  : build kernel image. It is the default sub-command.
- `modules`: build and install driver modules.
- `dtb`    : build DTB file.

## install

This command is used to invoke the install script of current profile. The script
resides in the directory of configuration file. 

If the script does not exist, a new script will be created and opened with the 
editor specified by `editor`. In this way, users can write an initial script.

The users can use command `edit install` to modify this script afterwards.

The arguments in the commandline will be transferred into the script. For example:

```
$ kbdashboard install arg1 arg2
```

arg1 and arg2 are parameters of the current installation script.

## make

This command is used to execute original targets of kernel in case existed
commands are not met the demands of users.

For example, the frist thing to do after getting the kernel source is often to 
make a default configuration. The command may be like:

```sh
$ make ARCH=arm bcm_defconfig
```

In such case, `bcm_defconfig` is the target argument:

```sh
$ kbdashboard make bcm_defconfig
```

# LICENSE
The GPLv3 License. See `LICENSE.md` file for more details.
