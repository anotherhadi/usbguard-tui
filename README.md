```
 _   _ ____  ____   ____                     _   _____ _   _ ___
| | | / ___|| __ ) / ___|_   _  __ _ _ __ __| | |_   _| | | |_ _|
| | | \___ \|  _ \| |  _| | | |/ _` | '__/ _` |   | | | | | || |
| |_| |___) | |_) | |_| | |_| | (_| | | | (_| |   | | | |_| || |
 \___/|____/|____/ \____|\__,_|\__,_|_|  \__,_|   |_|  \___/|___|
```

![GitHub Stars](https://www.shieldcn.dev/github/stars/anotherhadi/usbguard-tui.svg?variant=outline&theme=violet)
![Release](https://www.shieldcn.dev/github/release/anotherhadi/usbguard-tui.svg?variant=outline&theme=violet)
![CI](https://www.shieldcn.dev/github/ci/anotherhadi/usbguard-tui.svg?variant=outline&theme=violet)
[![Ko-fi](https://www.shieldcn.dev/badge/Ko--fi-sponsor-FF5E5B.svg?logo=kofi&variant=secondary&theme=violet)](https://ko-fi.com/anotherhadi)

# USBGuard TUI

> A terminal UI for managing USB devices via [usbguard](https://usbguard.github.io/).

USBGuard is a software framework for implementing a USB device authorization policy (allowlisting/blocklisting). It protects your system against rogue USB devices by scanning them and checking their parameters against a set of rules.

<img alt="USBGuard-tui demo" src="./.github/assets/demo.gif" width="600" />

Built with [bubbletea](https://github.com/charmbracelet/bubbletea) & Golang!

Colors and styles can be customized using [ilovetui](https://github.com/anotherhadi/ilovetui), which applies theme changes across all compatible TUI applications at once.

> This project is NOT affiliated with, endorsed by, or connected to "Red Hat, Inc" in any way. It is an unofficial, community-driven tool.

## Features

- List all connected USB devices with their current status
- Allow, block, or reject devices: temporarily or permanently
- Action popup for quick device management
- Filter devices by name with `/`
- Auto-refresh
- Keyboard shortcuts for all actions (`a`/`A`, `b`/`B`, `e`/`E`, ...) & Mouse support

## Requirements

- usbguard installed and the daemon running
- Sufficient privileges to communicate with the usbguard daemon socket

## Installation

<details>
<summary>Go install</summary>

```sh
go install github.com/anotherhadi/usbguard-tui@latest
```

</details>

<details>
<summary>Build from source</summary>

```sh
git clone https://github.com/anotherhadi/usbguard-tui
cd usbguard-tui
go build -o usbguard-tui ./cmd/usbguard-tui
```

</details>

<details>
<summary>NUR (Nix/NixOS)</summary>

Available via [NUR](https://github.com/nix-community/NUR), under the `anotherhadi` repo:

```sh
environment.systemPackages = [ nur.repos.anotherhadi.default-creds-tui ];
```

</details>

## Usage

```
usbguard-tui
```

The device list refreshes automatically every 2 seconds.

---

<div align="center">
  <a href="https://github.com/anotherhadi/usbguard-tui">github</a> |
  <a href="https://gitlab.com/anotherhadi_mirror/usbguard-tui">gitlab (mirror)</a> |
  <a href="https://git.hadi.icu/anotherhadi/usbguard-tui">gitea (mirror)</a>
</div>
