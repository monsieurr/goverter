## Goal of this project
Make a simpler easier to use converter.

## Project structure
.
├── README.md : the README file, you are here
├── main.go : GO Web server, backend stuff
├── package-lock.json : generate this with npm
├── package.json : generate this with npm
├── postcss.config.js : base postcss stuff (installed with tailwind)
├── src
│   └── input.css : to create my output.css file, should be put elsewhere probably
├── static
│   └── output.css : contains Tailwind css rules
├── tailwind.config.js : used to generate output.css
└── templates
    ├── index.html : main HTML frontend stuff
    └── result.html : deprecated / not used anymore

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