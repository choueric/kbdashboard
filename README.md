# Kernel Build Dashboard

[![Build Status](https://travis-ci.org/choueric/kbdashboard.svg?branch=master)](https://travis-ci.org/choueric/kbdashboard)

This tool (i.e. `kbdashboard`) is used to configure and manage building process
of multiple linux kernels. It is written in Golang.

It is just simply easy and comfortable to build ony one kernel with only one
configuration. But it is perfectly different when you have to handle with 
various kernels used in different projects or various configruations of one 
kernel. This tool helps you tackle with the management.

## features

- Run in any directories, no need of changing into the one where the kernel
  source tree is.
- Use individual building directory without affecting the kernel source tree.
- Simple commands to perform various actions from configuring to installing.
- Easy to configure by the json format.
- Colorful shell output.
- Built-in environment variables for installation scrips.
- Find out all DTS files related to the target DTB.


# Detailed Information

See [this post](http://ericnode.info/post/kbdashboard/).

# TODO
- [X] Built-in variables for installation scripts.
- [X] Add dts gathering.
- [X] Add extra options for kernel build, like `CFLAGS_KERNEL=-march=armv7-a`.
- [ ] Improve the install script template.
- [ ] Get the version string, include local version, like `3.14.28-132859-g953d55a`
- [ ] Add Dropbox support to sync and backup configurations and scripts.
- [ ] Complete test code.

# LICENSE
The GPLv3 License. See `LICENSE.md` file for more details.
