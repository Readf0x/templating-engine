{
  description = "What if we combined xml and golang?";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [ "x86_64-linux" ];
      perSystem = { system, pkgs, lib, ... }: {
        devShells.default = pkgs.mkShell {
          GOPATH = "/home/readf0x/.config/go";
          packages = [ pkgs.go pkgs.delve ];
        };
        packages = rec {
          templating-engine = pkgs.buildGoModule rec {
            name = "templating-engine";
            pname = "te";
            version = "release";

            src = ./.;

            buildInputs = [ pkgs.go ];

            vendorHash = null;

            postInstall = ''
              mkdir -p $out/share/man/man1
              mkdir -p $out/share/man/man7
              cp man/te.1 $out/share/man/man1
              cp man/te.7 $out/share/man/man7
              gzip $out/share/man/man1/te.1
              gzip $out/share/man/man7/te.7
            '';

            meta = {
              description = "The best templating engine";
              homepage = "https://github.com/Readf0x/templating-engine";
              license = lib.licenses.gpl3;
              mainProgram = pname;
            };
          };
          default = templating-engine;
        };
      };
    };
}
