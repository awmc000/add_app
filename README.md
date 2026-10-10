# Add App

A small Go + [Fyne](https://fyne.io) desktop utility for GNOME (Pop!_OS) that
collects launcher details in a form and writes valid `.desktop` files into the
user's applications directory (`~/.local/share/applications/`), so the app
appears in the app grid and dock.

## Requirements

- Go 1.22+
- GTK3 development headers (for Fyne): `libgtk-3-dev`

## Run

```sh
go run .
```

## Build

```sh
go build -o add-app .
```

## Usage

Fill in the form and press **Save .desktop**:

| Field | Maps to `.desktop` key | Notes |
|-------|------------------------|-------|
| Name | `Name=` | Also derives the filename (`my-app.desktop`) |
| Executable | `Exec=` | File picker; any path or command |
| Arguments | appended after executable | optional |
| Comment | `Comment=` | optional, tooltip/subtitle |
| Icon | `Icon=` | Browse… or Auto |
| Run in terminal | `Terminal=` | checkbox |
| Category | `Categories=` | dropdown; drives the default icon |

**Icon resolution** (highest first):
1. An explicitly picked icon.
2. An image beside the executable sharing its base name
   (`.png`, `.svg`, `.xpm`, `.jpg`, `.jpeg`, `.ico`).
3. A bundled default icon for the chosen category, copied to
   `~/.local/share/icons/`.

If a launcher with the same filename already exists, the app asks before
overwriting. Use **Load .desktop…** to open an existing file, edit it, and
save over it.

## Layout

```
internal/domain   pure logic (no Fyne imports), unit-tested
internal/app      Fyne UI: window, native file dialogs, controller wiring
assets/icons      bundled default category icons (SVG)
```

## Test

```sh
go test ./...
```

## Notes / current scope

- **User-only installs** (`~/.local/share/applications`). System-wide install
  (`/usr/share/applications`) is intentionally not implemented; a future
  enhancement could add it via polkit elevation.
- `Type=Application` and `Version=1.0` are hardcoded, since only app entries
  are created.
