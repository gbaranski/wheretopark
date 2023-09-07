const defaultTheme = require('tailwindcss/defaultTheme')

/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      animation: {
        ["infinite-slider"]: "infiniteSlider 20s linear infinite",
      },
      keyframes: {
        infiniteSlider: {
          "0%": { transform: "translateX(0)" },
          "100%": {
            transform: "translateX(calc(-250px * 5))",
          },
        },
      },
      fontFamily: {
        'sans': ['Poppins', ...defaultTheme.fontFamily.sans],
      }
    }
  },
  daisyui: {
      themes: [
        {
          myTheme: {
            "primary": "#FEE440",
            "secondary": "#383F51",
            "accent": "#EF2D56",
            "neutral": "#011627",
            "base-100": "#ffffff",
            "info": "#008BF8",
            "success": "#91F5AD",
            "warning": "#ec8004",
            "error": "#EF2D56",
          },
        },
      ],
  },
  plugins: [
    require("daisyui"), 
    require('@tailwindcss/typography'),
  ],
}

