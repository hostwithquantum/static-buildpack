# tests

## hugo example

```sh
hugo new site --format yaml ./hugo-example
```

```sh
hugo new site --format yaml ./hugo-npm
cd hugo-npm
npm create vite@latest . # setup with Vanilla javascript
npm install --save-dev tailwindcss @tailwindcss/cli @tailwindcss/vite
echo "/node_modules" >> .gitignore
rm -rf src/ dist/ # vite artefacts
```