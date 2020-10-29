<h1 align="center">
    Flippage v1.2.1
</h1>

![version](https://img.shields.io/badge/version-v1.2.1-brightgreen) ![README](https://img.shields.io/badge/README-available-brightgreen) ![platform](https://img.shields.io/badge/platform-mac-blue) ![maintainer](https://img.shields.io/badge/maintainer-Takuya%20HARA-brightgreen) 

A tiny app to flip book automatically. Just launch app, specify an interval to flip page in second and direction to flip. Then activate a book reader app and keep reading.

## How it works
Once the target app has selected Flippage sends key code periodically as if arrow key has pressed.

## Workaround
Unlike the others, Arrow Down is associated with KeyHold event as KeyDown is never invoked. Issue has opened here.
https://github.com/robotn/gohook/issues/26
