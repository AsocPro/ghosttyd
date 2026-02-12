![backend](https://github.com/tsl0922/ttyd/workflows/backend/badge.svg)
![frontend](https://github.com/tsl0922/ttyd/workflows/frontend/badge.svg)
[![GitHub Releases](https://img.shields.io/github/downloads/tsl0922/ttyd/total)](https://github.com/tsl0922/ttyd/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/tsl0922/ttyd)](https://hub.docker.com/r/tsl0922/ttyd)
[![Packaging status](https://repology.org/badge/tiny-repos/ttyd.svg)](https://repology.org/project/ttyd/versions)
![GitHub](https://img.shields.io/github/license/tsl0922/ttyd)

# ttyd - Share your terminal over the web

ttyd is a simple command-line tool for sharing terminal over the web.

![screenshot](https://github.com/tsl0922/ttyd/raw/main/screenshot.gif)

# Features

- Built on top of [libuv](https://libuv.org) and [ghostty-web](https://github.com/coder/ghostty-web) (WASM-compiled Ghostty terminal) for speed
- Fully-featured terminal with [CJK](https://en.wikipedia.org/wiki/CJK_characters) support and Unicode 15.1
- [ZMODEM](https://en.wikipedia.org/wiki/ZMODEM) ([lrzsz](https://ohse.de/uwe/software/lrzsz.html)) / [trzsz](https://trzsz.github.io) file transfer support
- Built-in URL detection and [OSC 8](https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda) hyperlink support
- SSL support based on [OpenSSL](https://www.openssl.org) / [Mbed TLS](https://github.com/Mbed-TLS/mbedtls)
- Run any custom command with options
- Basic authentication support and many other custom options
- Cross platform: macOS, Linux, FreeBSD/OpenBSD, [OpenWrt](https://openwrt.org), Windows

> ❤ Special thanks to [JetBrains](https://www.jetbrains.com/?from=ttyd) for sponsoring the opensource license to this project.

# Installation

## Install on macOS

- Install with [Homebrew](http://brew.sh): `brew install ttyd`
- Install with [MacPorts](https://www.macports.org): `sudo port install ttyd`

## Install on Linux

- Binary version (recommended): download from the [releases](https://github.com/tsl0922/ttyd/releases) page
- Install with [Homebrew](https://docs.brew.sh/Homebrew-on-Linux) : `brew install ttyd`
- Install the snap: `sudo snap install ttyd --classic`
- Build from source (debian/ubuntu):
    ```bash
    sudo apt-get update
    sudo apt-get install -y build-essential cmake git libjson-c-dev libwebsockets-dev
    git clone https://github.com/tsl0922/ttyd.git
    cd ttyd && mkdir build && cd build
    cmake ..
    make && sudo make install
    ```
    You may also need to compile/install [libwebsockets](https://libwebsockets.org) from source if the `libwebsockets-dev` package is outdated.
- Install on OpenWrt: `opkg install ttyd`
- Install on Gentoo: clone the [repo](https://bitbucket.org/mgpagano/ttyd/src/master) and follow the directions [here](https://wiki.gentoo.org/wiki/Custom_repository#Creating_a_local_repository).

## Install on Windows

- Binary version (recommended): download from the [releases](https://github.com/tsl0922/ttyd/releases) page
- Install with [WinGet](https://github.com/microsoft/winget-cli): `winget install tsl0922.ttyd`
- Install with [Scoop](https://scoop.sh/#/apps?q=ttyd&s=2&d=1&o=true): `scoop install ttyd`
- [Compile on Windows](https://github.com/tsl0922/ttyd/wiki/Compile-on-Windows)

## Building with Dagger

The project includes a [Dagger](https://dagger.io) pipeline that builds everything inside containers -- no need to install Node.js, Yarn, cmake, or any C libraries on your host. The only prerequisite is the [Dagger CLI](https://docs.dagger.io/install/).

**Full build** (frontend + C binary):

```bash
dagger call build --source=. export --path=./build-out
# produces build-out/ttyd
```

This runs the complete pipeline: installs JS dependencies, builds the ghostty-web frontend with webpack, generates the C header files (`html.h` and `wasm.h`), then compiles the ttyd binary with cmake.

**Regenerate frontend headers** (after updating `ghostty-web` or changing frontend code):

```bash
dagger call generate --source=. export --path=./src
```

This rebuilds just the frontend and writes the updated `src/html.h` and `src/wasm.h`. Useful when you want to iterate on the C code locally with `cmake --build build` without re-running the full pipeline.

**Update yarn.lock** (after changing `package.json`):

```bash
dagger call yarn-install --source=. export --path=./html
```

All available pipeline functions:

| Function | Description |
|---|---|
| `build` | Full build: frontend + C binary |
| `generate` | Regenerate `src/html.h` and `src/wasm.h` from frontend source |
| `frontend` | Same as `generate` (returns the header files as a directory) |
| `yarn-install` | Run `yarn install`, exports updated `html/` with new `yarn.lock` |
| `build-local` | Like `build`, but takes an `--output` path directly |

# Usage

## Command-line Options

```
USAGE:
    ttyd [options] <command> [<arguments...>]

OPTIONS:
    -p, --port              Port to listen (default: 7681, use `0` for random port)
    -i, --interface         Network interface to bind (eg: eth0), or UNIX domain socket path (eg: /var/run/ttyd.sock)
    -U, --socket-owner      User owner of the UNIX domain socket file, when enabled (eg: user:group)
    -c, --credential        Credential for basic authentication (format: username:password)
    -H, --auth-header       HTTP Header name for auth proxy, this will configure ttyd to let a HTTP reverse proxy handle authentication
    -u, --uid               User id to run with
    -g, --gid               Group id to run with
    -s, --signal            Signal to send to the command when exit it (default: 1, SIGHUP)
    -w, --cwd               Working directory to be set for the child program
    -a, --url-arg           Allow client to send command line arguments in URL (eg: http://localhost:7681?arg=foo&arg=bar)
    -W, --writable          Allow clients to write to the TTY (readonly by default)
    -t, --client-option     Send option to client (format: key=value), repeat to add more options
    -T, --terminal-type     Terminal type to report, default: xterm-256color
    -O, --check-origin      Do not allow websocket connection from different origin
    -m, --max-clients       Maximum clients to support (default: 0, no limit)
    -o, --once              Accept only one client and exit on disconnection
    -q, --exit-no-conn      Exit on all clients disconnection
    -B, --browser           Open terminal with the default system browser
    -I, --index             Custom index.html path
    -b, --base-path         Expected base path for requests coming from a reverse proxy (eg: /mounted/here, max length: 128)
    -P, --ping-interval     Websocket ping interval(sec) (default: 5)
    -6, --ipv6              Enable IPv6 support
    -S, --ssl               Enable SSL
    -C, --ssl-cert          SSL certificate file path
    -K, --ssl-key           SSL key file path
    -A, --ssl-ca            SSL CA file path for client certificate verification
    -d, --debug             Set log level (default: 7)
    -v, --version           Print the version and exit
    -h, --help              Print this text and exit
```

Read the example usage on the [wiki](https://github.com/tsl0922/ttyd/wiki/Example-Usage).

## Browser Support

Modern browsers with WebAssembly support (Chrome, Firefox, Safari, Edge).

## Alternatives

* [Wetty](https://github.com/krishnasrinivas/wetty): [Node](https://nodejs.org) based web terminal (SSH/login)
* [GoTTY](https://github.com/yudai/gotty): [Go](https://golang.org) based web terminal
