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
              $out/bin/te man/te.7.tet
              for i in {1,7}; do
                mkdir -p $out/share/man/man$i
                cp man/te.$i $out/share/man/man$i
              done
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
