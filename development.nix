{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  nativeBuildInputs = [
    pkgs.pkg-config
    pkgs.tesseract
    pkgs.leptonica
    pkgs.glib
  ];

  shellHook = ''
    export C_INCLUDE_PATH="${pkgs.leptonica}/include:${pkgs.glib.dev}/include/glib-2.0:${pkgs.glib.out}/lib/glib-2.0/include:$C_INCLUDE_PATH"
    export LIBRARY_PATH="${pkgs.leptonica}/lib:${pkgs.glib.out}/lib:$LIBRARY_PATH"
  '';
}
