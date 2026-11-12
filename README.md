```
 █ █ █▀▀ █▄▄ █▀▀ █ █ ▄▀█ █▀█ █▀▄ ▄▄ ▀█▀ █ █ █
 █▄█ ▄▄█ █▄█ █▄█ █▄█ █▀█ █▀▄ █▄▀     █  █▄█ █
```

![GitHub Stars](https://www.shieldcn.dev/github/stars/anotherhadi/usbguard-tui.svg?variant=outline&theme=violet)
![Release](https://www.shieldcn.dev/github/release/anotherhadi/usbguard-tui.svg?variant=outline&theme=violet)
![CI](https://www.shieldcn.dev/github/ci/anotherhadi/usbguard-tui.svg?variant=outline&theme=violet)

# USBGuard TUI

> A terminal UI for managing USB devices via [usbguard](https://usbguard.github.io/).

USBGuard is a software framework for implementing a USB device authorization policy (allowlisting/blocklisting). It protects your system against rogue USB devices by scanning them and checking their parameters against a set of rules.

<img alt="USBGuard-tui demo" src="./.github/assets/demo.gif" width="600" />

Built with [bubbletea](https://github.com/charmbracelet/bubbletea) & Golang!

Built on [foam](https://github.com/anotherhadi/foam), which shares its theme, config and keybinds across compatible TUI applications.

> This project is NOT affiliated with, endorsed by, or connected to "Red Hat, Inc" in any way. It is an unofficial, community-driven tool.

## Features

- List all connected USB devices with their current status
- Allow, block, or reject devices: temporarily or permanently
- Action overlay for quick device management
- Filter devices by name with `/`
- Command palette (`ctrl+p`) and a help overlay (`?`) listing every shortcut
- Auto-refresh
- Keyboard shortcuts for all actions (`a`/`A`, `b`/`B`, `e`/`E`, ...) & Mouse support
- Configurable theme and keybinds, shared with other [foam](https://github.com/anotherhadi/foam) tools

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
environment.systemPackages = [ nur.repos.anotherhadi.usbguard-tui ];
```

</details>

## Usage

```
usbguard-tui
```

The device list refreshes automatically every 2 seconds.

## Configuration

To change the defaults (theme, keybinds, refresh interval), put a
`~/.config/usbguard-tui/config.yaml` file; you only need the keys you
want to change, the rest falls back to the defaults. Theme and keybinds
shared with other foam tools can also go in `~/.config/foam/config.yaml`.
Run `usbguard-tui --print-config` to see the merged config, and see
[foam's defaults](https://github.com/anotherhadi/foam/blob/main/defaults.yaml)
for the theme and the keys every foam tool shares.

---

<div align="center">
  <a href="https://github.com/anotherhadi/usbguard-tui">github</a> |
  <a href="https://git.hadi.icu/anotherhadi/usbguard-tui">gitea (mirror)</a>
</div>
