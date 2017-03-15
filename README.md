# Kernel Build Dashboard

[![Build Status](https://travis-ci.org/choueric/kbdashboard.svg?branch=master)](https://travis-ci.org/choueric/kbdashboard)

This tool (i.e. `kbdashboard`) is used to configure and manage building process
of multiple linux kernels. It is written in Golang.

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
- [ ] Add Dropbox support to sync and backup configurations and scripts.
- [ ] Complete test code.

# LICENSE
The GPLv3 License. See `LICENSE.md` file for more details.
