/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./ui/*.html"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    themes: ["cupcake"],
  },
};
