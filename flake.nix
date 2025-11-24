{
  description = "NC CMS";

	inputs = {
		nixpkgs.url = "github:NixOS/nixpkgs";
		flake-utils.url = "github:numtide/flake-utils";
		gomod2nix.url = "github:nix-community/gomod2nix";
		gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
		gomod2nix.inputs.flake-utils.follows = "flake-utils";
	};

  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
		flake-utils.lib.eachDefaultSystem (system:
			let
				pkgs = import nixpkgs {
					inherit system;
					overlays = [ gomod2nix.overlays.default ];
				};

				php-dev = (pkgs.php.overrideAttrs (old: {
					configureFlags = old.configureFlags ++ [
						"--enable-embed"
					];
				})).unwrapped.dev;
				php-dev-cflags = "-I${php-dev}/include/php -I${php-dev}/include/php/main -I${php-dev}/include/php/TSRM -I${php-dev}/include/php/Zend -I${php-dev}/include/php/ext -I${php-dev}/include/php/ext/date/lib";
				php-dev-ldflags = "-L${php-dev}/lib -lphp -Wl,-rpath,${php-dev}";
			in {
				packages.default = pkgs.buildGoApplication rec {
					pname = "nccms";
					version = "0.1.0";
					src = ./.;
					modules = ./gomod2nix.toml;
					vendorHash = null;
					buildInputs = with pkgs; [
						gcc
						php-dev
						pkg-config
					];
					CGO_CFLAGS = php-dev-cflags;
					CGO_LDFLAGS = php-dev-ldflags;
				};

				apps.default = {
					type = "app";
					program = "${self.packages.${system}.default}/bin/nccms";
				};

				devShells.default = pkgs.mkShell {
					shellHook = ''
						export PS1="\[\033[0;1;32m\][nccms dev]$ \[\033[0m\]"

						export CGO_CFLAGS="${php-dev-cflags}"
						export CGO_LDFLAGS="${php-dev-ldflags}"
					'';
					buildInputs = with pkgs; [
						php-dev
						pkg-config
					];
				};
			});
}
