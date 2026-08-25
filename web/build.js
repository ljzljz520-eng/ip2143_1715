const fs = require('fs');
const path = require('path');
const source = fs.readFileSync(path.join(__dirname, 'index.html'), 'utf8');
fs.mkdirSync(path.join(__dirname, 'dist'), { recursive: true });
fs.writeFileSync(path.join(__dirname, 'dist', 'index.html'), source.replace('<!-- generated -->', '<!-- generated -->'));
