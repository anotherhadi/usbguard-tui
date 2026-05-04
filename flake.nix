{
  description = "A terminal UI for managing USB devices via usbguard.";

  inputs = {nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";};

  outputs = {
    self,
    nixpkgs,
  }: let
    supportedSystems = ["x86_64-linux" "aarch64-linux"];

    forAllSystems = f:
      nixpkgs.lib.genAttrs supportedSystems
      (system: f system (import nixpkgs {inherit system;}));

    pname = "usbguard-tui";
    version = "1.0.1";

    ldflags = ["-s" "-w" "-X main.version=${version}"];
  in {
    packages = forAllSystems (system: pkgs: let
      pkg = pkgs.buildGoModule {
        inherit pname version ldflags;

        src = ./.;
        outputs = ["out"];

        vendorHash = "sha256-NfmNdQaISob3ZguQnwgfXHUDUFION988ucp18q996pM=";

        meta = with pkgs.lib; {
          description = "A terminal UI for managing USB devices via usbguard.";
          homepage = "https://github.com/anotherhadi/usbguard-tui";
          platforms = platforms.unix;
        };
      };
    in {
      "${pname}" = pkg;
      default = pkg;
    });
  };
}
