# HLive-Wails Hacker News Reader

A Hacker News reader built using Wails and HLive.

This project is currently a prototype to create a development process for HLive developers to use Wails. Eventually, 
the goal is to turn this into an example app for this tech.

## Features

- Custom build process
  - Tailwind CSS
  - esbuild JavaScript processing
  - Copy images, fonts, etc.
- Custom Dev process
  - Go file server for static assets
  - Known issue: need to reload browser on rebuild
- Custom CSS
  - frontend/css/app.css
- Custom JavaScript
  - frontend/js/app.js

## Notes

- The Hacker News API returns raw HTML. We use JavaScript to catch clicks on them and open the system browser.
- The index.html file is a simple redirect to an HLive server.

## TODO

- Move dev and page server to an official repo

### Bugs

- Reload browser on rebuild when in dev mode

### Features

- Reload static content on change
- Add custom fonts
- Add support for user setting
  - Store on disk
- Cache data to disk
- Reload button
- Light and dark modes
- Button to open comments in browser
- Comment navigation
  - next, prev, parent, root, next root
- Remember story clicked state
- Add other lists in addition to "Top Stories"
  - Best, New, Ask,, etc.
- Use social media style preview
  - https://socialsharepreview.com/
- Show fav icon
- Show domain
- Scroll back to top
- Infinite scroll
- Add dev server URL config
- File upload
  - This will be reworked in HLive soon so wait for that first.