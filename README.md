# AirMonGo - Wi-Fi Scanner

**AirMonGo** is a Wi-Fi network scanner for Windows, Linux, and macOS that lists available networks, their signal strength, security type, and more.

## Features
- Supports Windows, Linux, and macOS
- Lists available Wi-Fi networks with details like signal strength, encryption, authentication, and channel
- Provides signal quality (Excellent, Good, Fair, Poor)
- macOS 14.4+ compatible with modern Wi-Fi scanning commands

## Installation

### Windows:
1. Make sure you have the necessary tools (`netsh`)
2. Download or clone this repo
3. Run the tool via the command line

### Linux:
1. Install `nmcli` if not already installed
2. Download or clone this repo
3. Run the tool via the command line

### macOS:
1. Requires `wdutil` (macOS 14.4+) or `system_profiler` (fallback for older versions)
2. Download or clone this repo
3. Run the tool via the command line

## Usage

Run the program from the command line:
```bash
AirMonGo
```

For help:
```bash
AirMonGo --help
```

## Compatibility

- **Windows**: Windows 7 and later
- **Linux**: Any distribution with NetworkManager and `nmcli`
- **macOS**: macOS 10.13+ (uses `wdutil` on 14.4+, `system_profiler` on older versions)

## Version History

### v1.2
- Fixed macOS 14.4+ compatibility by replacing deprecated `airport` command
- Added `wdutil` support with `system_profiler` fallback
- Improved signal strength handling for dBm values

### v1.1
- Initial release

## Contributing

Thanks to:
- [@brianhenrydev](https://github.com/brianhenrydev) - Reported macOS 14.4+ `airport` deprecation issue ([#1](https://github.com/staxsum/AirMonGo/issues/1))

## License

MIT
