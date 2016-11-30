# Kernel Build Dashboard
This tool, i.e. `kbdashboard`, is a dashboard used to configure and manage build
process of multiple linux kernels. It is written in Golang.

Developpers, especially in embedded system, may usually need to modify, build 
and test more than one linux kernels. Some may handle different kernel versions
using different toolchains, some may handle one kernel version but use different
build directories based on different configuratios. In these cases, users have
to remember various configurations and use all kinds of commands.

In order to simplify this matter, `kbdashboard` acts as the wrapper to handle
all these complicated details. It has some advantages as below:

- Run in any directory without changing to the direcotry of kernel sources.
- Use individual building directory without affecting the kernel source tree.
- Simple commands to perform various actions from starting up to installing.
- Configure easily by using the json format configuration file.
- Colorful shell output.

# Build this Tool
```
$ cd kernelBuildDashboard
$ make
$ sudo make install
```

Then user will get the executable `kbdashboard`. It will be installed into the
direcotry `/usr/local/bin`. And the bash completion script will be installed
into `/etc/bash_completion.d`. The Makefile is very simple and can be modified
easily.

# How to Use this Tool

## Create a Profile

The configuration file of `kbdashboard` is the core. Use command:
```sh
$ kbdashboard edit
```
to edit the json format configuration file. If it is the first time to execute
this tool, there is no configuration file `$HOME/.config/kbdashboard/config.json`.
The tool will create a template and open it using `vim`. The template is:

```json
{
	"editor": "vim",
	"current": 0,
	"profile": [
	{
		"name":"demo",
		"src_dir":"/home/user/kernel"
		"arch":"arm",
		"cross_compile":"arm-eabi-",
		"target":"uImage",
		"output_dir":"./_build",
		"defconfig":"at91rm9200_defconfig",
		"dtb":"at91rm9200ek.dtb",
		"mod_install_dir":"./_build/mod",
		"thread_num":4,
	}
	]
}
```

There are some global options and only one profile with the name `demo`.

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

Others can be empty and will be the default value during building kernel.

```
arch            : architecture, corresponding to `ARCH` of kernel build command.
cross_compile   : cross compiler, corresponding to `CROSS_COMPILE` of kernel 
                  build command.
target          : target of building kernel image.
output_dir      : output build directory, corresponding to `O` of kernel build 
                  command. It is relative to the `src_dir` if not absolute.
defconfig       : default configuration when start up.
dtb             : the name of target DTB file.
mod_install_dir : module install directory, corresponding to `INSTALL_MOD_PATH`
                  of kernel build command. It is relative to the `src_dir` if
				  not absolute.
thread_num      : number of thread used to compile, corresponding to `-j` option.
```

After the editing, use command `list` to see the result:

```sh
$ kbdashboard list
```

Since `demo` is listed, use command `choose` to set it as the default
profile:

```sh
$ kbdashboard choose demo
```

Use `list` command to see whether `demo` is marked by red asterisk.
From now on, it is no need to append profile's name or index in commands.

## Configure Kernel

Before comiling the kernel, it has to do the default configuration:

```sh
$ kbdashboard config def
```

It will use `defconfig` in profile to generate `.config` file in the buiding
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

# Details of Tool Commands

## help
Use command `help` or do not use any command to show the help message:

```sh
$ kbdashboard help
or
$ kbdashboard
```

`help` can also be used as sub-command in some commands. For example:

```sh
$ kbdashboard build help
```

## list
List the profiles. The current profile is marked by the red asterisk symbol.

sub-command:

- `verbose`: show all information of profiles.

## choose
Choose the current profile by name or index.

## edit
Edit the configuration file using editor specified by the `editor` configuration.

sub-command:

- `config` : edit the tool's configuration file. It is the default sub-command.
- `install`: edit one profile's installation script.

## config
Operate the configurations. The profile is specified by the name or index of
profile in the command line, or by the current index in the configuration file
if no option in the command line.

sub-command:

- `menu`: invoke `menuconfig`. It is the default sub-command.
- `def` : invoke `xxxx_defcofig` which specified by `defcofig` of the profile 
          configuration.
- `save`: invoke `savedefconfig`.

## build
Build the target for specified kernel profile. 

- `image`  : build kernel image. It is the default sub-command.
- `modules`: build and install driver modules.
- `dtb`    : build DTB file.

## install
This command is used to call the install script of specified profile. The script
is in the directory of configuration file. 

If there is no such script, a new script will be created and opened by the 
editor which is specified by `editor` config. In this way, users can write a 
initial script.

The users can also use the command `edit install` to call the editor and modify
this script explicitly.

## make
This command is used to execute original targets of kernel.

For example, the frist thing to do after getting the kernel source is often to 
make a default configuration. The command may be like:

```sh
$ make ARCH=arm bcm_defconfig
```

In such case, `bcm_defconfig` is the target argument:

```sh
$ kbdashboard make bcm_defconfig
```

The command uses the chosen profile. Or to specify the first profile:

```sh
$ kbdashboard make bcm_defconfig 0
```

Another useful targets are `modules`, `install` and `modules_install`. 

`modules` is used to compile kernel modules.

`install` is one target belonged to kernel's make system, and it's different 
from the command `install`. The target invokes an install sript for specified 
architecture. For example, building a kernel for a Debian host. In this case, 
the architecture is x86, script `arch/x86/boot/install.sh` is invoked.

`modules_install` installs driver modules to a directory whose path is
`$(mod_install_dir)/lib/modules/$(KERNELRELEASE)`. The path would be in
`/lib/modules` if the `mod_install_dir` is empty in configuration, and it is 
the common path for building a kernel for PC.


# LICENSE
The GPLv3 License. See `LICENSE.md` file for more details.
