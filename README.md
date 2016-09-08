# kernel Build Dashboard
Dashboard for configuring, managing build process of multiple linux kernel.

# build
```
$ cd kbd_cli
$ go build
```

the result executable is `kbd_cli`.

# configuration
Configuration file is in json format. It can contains multiple kernel
configurations; each one is a profile. The program finds configuration file 
`~/.config/kbdashboard/config.json`.

A sample is shown below:

```json
{
	"profile": [
	{
		"name":"demo",
		"src_dir":"/home/user/kernel"
		"arch":"arm",
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

One profile must include following values:
```
name    : profile name.
src_dir : directory path of kernel source.
```

Values below are optional:
```
arch            : architecture, corresponding to `ARCH` of kernel build command.
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
