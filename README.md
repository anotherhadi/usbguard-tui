```
▖▖▄▖▄ ▄▖       ▌  ▄▖▖▖▄▖
▌▌▚ ▙▘▌ ▌▌▀▌▛▘▛▌  ▐ ▌▌▐
▙▌▄▌▙▘▙▌▙▌█▌▌ ▙▌  ▐ ▙▌▟▖

```

# USBGuard TUI

A terminal UI for managing USB devices via [usbguard](https://usbguard.github.io/).

USBGuard is a software framework for implementing a USB device authorization policy (allowlisting/blocklisting). It protects your system against rogue USB devices by scanning them and checking their parameters against a set of rules.

Built with [bubbletea](https://github.com/charmbracelet/bubbletea) & Goland!

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
go build -o usbguard-tui .
```

</details>

<details>
<summary>Nix run</summary>

```sh
nix run github:anotherhadi/usbguard-tui
```

</details>

<details>
<summary>NixOS: system installation</summary>

Add the flake input and include the package in your configuration:

```nix
# flake.nix
inputs.usbguard-tui.url = "github:anotherhadi/usbguard-tui";

# configuration.nix / home.nix
environment.systemPackages = [ inputs.usbguard-tui.packages.${system}.default ];
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
</div
