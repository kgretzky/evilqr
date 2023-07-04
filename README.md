# Evil QR

Toolkit demonstrating another approach of a QRLJacking attack, allowing to perform remote account takeover, through sign-in QR code phishing.

It consists of a browser extension used by the attacker to extract the sign-in QR code and a server application, which retrieves the sign-in QR codes to display them on the hosted phishing pages.

Read more about it on my blog: https://breakdev.org/evil-qr-phishing

## Configuration

The parameters used by **Evil QR** are hardcoded into extension and server source code, so it is important to change them to use custom values, before you build and deploy the toolkit.

| parameter | description | default value |
|-|-|-|
| **API_TOKEN** | API token used to authenticate with REST API endpoints hosted on the server | 00000000-0000-0000-0000-000000000000 |
| **QRCODE_ID** | QR code ID used to bind the extracted QR code with the one displayed on the phishing page | 11111111-1111-1111-1111-111111111111 |
| **BIND_ADDRESS** | IP address with port the HTTP server will be listening on | 127.0.0.1:35000 |
| **API_URL** | External URL pointing to the server, where the phishing page will be hosted | http://127.0.0.1:35000 |

Here are all the places in the source code, where the values should be modified:

#### server/core/config.go:
```
const API_TOKEN = "00000000-0000-0000-0000-000000000000"
const BIND_ADDRESS = "127.0.0.1:35000"
```

#### server/templates/index.html:
```
const API_URL = "http://127.0.0.1:35000";
const QRCODE_ID = "11111111-1111-1111-1111-111111111111";
```

#### extension/background.js:
```
const API_TOKEN = "00000000-0000-0000-0000-000000000000";
const API_URL = "http://127.0.0.1:35000";
const QRCODE_ID = "11111111-1111-1111-1111-111111111111";
```

## Installation

### Extension

You can load the extension in Chrome, through `Load unpacked` feature:
https://developer.chrome.com/docs/extensions/mv3/getstarted/development-basics/#load-unpacked

Once the extension is installed, make sure to pin its icon in Chrome's extension toolbar, so that the icon is always visible.

### Server

Make sure you have [Go installed](https://go.dev/doc/install) version at least 1.20.

To build go to `/server` directory and run the command:

Windows:
```
build_run.bat
```

Linux:
```
chmod 700 build.sh
./build.sh
```

Built server binaries will be placed in the `./build/` directory.

## Usage

1. Run the server by running the built server binary: `./server/build/evilqr-server`
2. Open any of the supported websites in your Chrome browser, with installed **Evil QR** extension:
```
https://discord.com/login
https://web.telegram.org/k/
https://whatsapp.com
https://store.steampowered.com/login/
https://accounts.binance.com/en/login
https://www.tiktok.com/login
```
3. Make sure the sign-in QR code is visible and click the **Evil QR** extension icon in the toolbar. If the QR code is recognized, the icon should light up with colors.
4. Open the server's phishing page URL: `http://127.0.0.1:35000` (default)

## License

**Evil QR** is made by Kuba Gretzky ([@mrgretzky](https://twitter.com/mrgretzky)) and it's released under MIT license.
