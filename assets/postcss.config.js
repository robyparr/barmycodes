module.exports = {
  plugins: [
    require('postcss-import'),
    require('postcss-flexbugs-fixes'),
    require('postcss-nested'),
    require('tailwindcss')('./js/tailwindcss/tailwind.config.js'),
    require('autoprefixer'),
  ]
}
