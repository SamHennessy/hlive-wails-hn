const fs = require('fs');

fs.rmSync("./dist", { recursive: true, force: true });

require('child_process').execSync(
    'npx tailwindcss -i ./css/app.css -o ./dist/css/app.css',
    {stdio: 'inherit'}
);

require('esbuild').build({
    entryPoints: ['./js/app.js'],
    bundle: true,
    minify: true,
    sourcemap: true,
    target: ['chrome58', 'firefox57', 'safari11', 'edge16'],
    outfile: './dist/js/app.js',
}).catch(() => process.exit(1));

const copyfiles = require("copyfiles");

copyfiles(["index.html","./fonts/**", "./img/**", "./dist/"], {}, (err) => {
    if (err) {
        console.log("Error occurred while copying",err);
    }
    console.log("folder(s) copied to destination");
});

