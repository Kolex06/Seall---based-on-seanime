# Seall

Seall is a self-hosted media server and desktop/web app being reshaped into an all-media library for movies, TV shows, anime, and local files.

The main data provider is SIMKL. SIMKL is used because it tracks TV, movies, and anime in one account and exposes sync endpoints for all of them.

## Current Direction

- SIMKL-first authentication and media sync.
- Local library scanning and playback for all media.
- Movie, TV show, and anime entries flowing through one media collection.
- Legacy SIMKL-shaped internal types are being kept temporarily so playback, scanning, and UI modules can be migrated safely instead of all at once.
- Manga-specific features are not part of the SIMKL replacement path unless a separate provider is added later.

## SIMKL Setup

Create a SIMKL application at <https://simkl.com/settings/developer/> and set the client ID in one of these places:

- Environment variable: `SIMKL_CLIENT_ID`
- Config file: `simkl.clientId`

For OAuth code exchange, also set:

- `SIMKL_CLIENT_SECRET`
- `SIMKL_REDIRECT_URI`

The app also includes a SIMKL PIN login path for local/desktop usage.

## Development

The Go module path is `seall`.

## Marketplace

The default Seall extension Marketplace is stored in this repository:

- `marketplace.json` is used by the in-app Marketplace.
- `extensions.json` is a repository import list for bulk extension installs.
