/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./ui/*.html"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["cupcake"],
  },
};
