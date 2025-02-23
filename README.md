# README

## About

This is a Wails go-app template.

It combines the wails with [go-app](https://github.com/maxence-charriere/go-app) to build the front-end in pure go and
compile it to WebAssembly (WASM).

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: <https://wails.io/docs/reference/project-config>

## Live Development

To run in live development mode, run `wails dev` in the project directory. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on <http://localhost:34115>. Connect
to this in your browser, and you can call your Go code from DevTools.

## Building

To build a redistributable, production mode package, use `wails build`.
