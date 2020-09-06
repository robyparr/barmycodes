// We need to import the CSS so that webpack will load it.
// The MiniCssExtractPlugin is used to separate it out into
// its own CSS file.
import "./tailwindcss/application.css"

// webpack automatically bundles all modules in your
// entry points. Those entry points can be configured
// in "webpack.config.js".
//
// Import deps with the dep name or local files with a relative path, for example:
//
//     import {Socket} from "phoenix"
//     import socket from "./socket"
//
import "phoenix_html"
import Sidebar from "./components/Sidebar"
window.Sidebar = Sidebar

const feather = require("feather-icons")
window.renderIcons =
  function() {
    document.querySelectorAll('.select-wrapper').forEach(wrapper => {
      var iconChild = document.createElement('i')
      iconChild.dataset.feather = 'chevron-down'
      iconChild.classList.add('chevron')

      wrapper.appendChild(iconChild)
    })

    document.querySelectorAll('[data-feather=""]').forEach(el => el.remove())
    feather.replace()
  }
document.addEventListener('DOMContentLoaded', renderIcons)
