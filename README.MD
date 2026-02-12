# FFConvert

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org/)

FFConvert is a lightweight system-level file converter that automatically converts downloads when you change the file extension.

**Example:**

- Download `image.webp`
- Rename it to `image.png` in the save dialog
- FFConvert automatically converts it to a real PNG

No websites. No manual converters. No waiting.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Uninstall](#uninstall)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Automatic conversion on download**: Works directly from browser save dialog
- **Image conversions**:
  - webp ↔ png ↔ jpg
  - avif → png/jpg
- **Video container conversions**:
  - mp4 ↔ mkv
  - webm → mp4/mkv
- Uses FFmpeg for reliable conversions
- Native Windows integration

## Requirements

- Windows
- [Go](https://golang.org/dl/) installed
- [FFmpeg](https://ffmpeg.org/download.html) installed and in PATH

Test FFmpeg:

```bash
ffmpeg -version
```

## Installation

1. Clone or download this repository.
2. Run:

   ```bash
   install.bat
   ```

3. Open browser extensions page:

   - Brave: `brave://extensions`
   - Chrome: `chrome://extensions`

4. Enable **Developer mode**.
5. Click **Load unpacked**.
6. Select the `extension` folder.
7. Copy the extension ID.
8. Open:

   ```
   C:\ffconvert\ffconvert-helper.json
   ```

9. Replace:

   ```
   REPLACE_EXTENSION_ID
   ```

   with your extension ID.

10. Run:

    ```bash
    install.bat
    ```

    again.

Installation complete.

## Usage

1. Download a file.
2. Change the extension in the save dialog.
3. File converts automatically.

**Example:**

```
video.mp4 → save as video.mkv
```

## Uninstall

Run:

```bash
uninstall.bat
```

Then remove the extension from your browser.

## Roadmap

### v2
- Document conversions
  - docx → txt
  - pdf → text
- Audio conversions
- Conversion rules

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
