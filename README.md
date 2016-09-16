# kernel Build Dashboard
Dashboard for configuring, managing build process of multiple linux kernel.

# build
```
$ cd kernelBuildDashboard
$ make
```

the result executable is `kbdashboard`.

# configuration
Configuration file is in json format. It can contains multiple kernel
configurations; each one is a profile. The program finds configuration file 
`~/.config/kbdashboard/config.json`.

A sample is shown below:

```json
{
	"editor": "vim",
	"current": 0,
	"profile": [
	{
		"name":"demo",
		"src_dir":"/home/user/kernel"
		"arch":"arm",
		"target":"uImage",
		"cross_compile":"arm-eabi-",
		"output_dir":"./_build",
		"mod_install_dir":"./_build/mod",
		"thread_num":4,
	},
	{
		"name":"demo2",
		"src_dir":"/home/user/kernel2"
	}
	]
}
```

The options below are globel:
```
editor  : Specify text editor which will be invoked when 'edit' command is executed.
current : Current profile index. If no speficy profile in 'make' and 'config' command,
          this index profile will be used.
``` 

One profile must include following values:
```
name    : profile name.
src_dir : directory path of kernel source.
```

Values below are optional:
```
arch            : architecture, corresponding to `ARCH` of kernel build command.
target          : target of the make command.
cross_compile   : cross compiler, corresponding to `CROSS_COMPILE` of kernel 
                  build command.
output_dir      : output build directory, corresponding to `O` of kernel build 
                  command.
mod_install_dir : module install directory, corresponding to `INSTALL_MODE_PATH`
                  of kernel build command.
thread_num      : number of thread used to compile, corresponding to `-j` option.
```

If these options are not specified in configuration file, programe just ignores
them.

# Commands

## help
Now there are 6 commands which are shown via command 'help' command:

```
$ kbdashboard help
Usage:
  - edit        : Edit the config file using editor specified in config file.
  - config      : [name | index]. Configure kernel using menuconfig
  - choose      : {name | index}. Choose current profile.
  - list        : List all profiles
  - make        : [name | index]. Build kernel specified by name or index
  - help        : Display this message.
```

## list
List the profiles. The current profile is marked by the red '*'.

## edit
Edit the configuration file using editor specified by the "editor" option.

## config
Invoke menuconfig to the specified kernel profile. The profile is specified by
the name or index of profile in the command line, or by the current index in
the configuration file if no option in the command line.

## make
Make the target for specified kernel profile. The way to specify profile is as
same as command 'config'.

## choose
Choose the current profile by name or index.

