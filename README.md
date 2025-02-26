## Goal of this project
Make a simpler easier to use converter.

## Stack
- Go
- Little bit of HTMX
- CSS (using Tailwind)

## Commands
```bash
npx tailwindcss -i ./src/input.css -o ./static/output.css --minify
npm run build # Same as the previous line but shorter
go run main.go # Launch the local server (port 8080)
```

## Current features
- Converts common units
- Copy results
- Dark mode toggle

## Potential future updates
- Adding more units
- History
- Deployment somewhere