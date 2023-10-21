import defaultTheme from 'tailwindcss/defaultTheme';

const theme = {
  "primary": "#FEE440",
  "secondary": "#383F51",
  "accent": "#EF2D56",
  "neutral": "#011627",
  "base-100": "#ffffff",
  "info": "#008BF8",
  "success": "#91F5AD",
  "warning": "#ec8004",
  "error": "#EF2D56",
};


/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      fontFamily: {
        'sans': ['Montserrat', ...defaultTheme.fontFamily.sans],
      },
      keyframes: {
        typing: {
          "0%": {
            width: "0%",
            visibility: "hidden"
          },
          "100%": {
            width: "100%"
          }  
        },
        blink: {
          "50%": {
            borderColor: "transparent"
          },
          "100%": {
            borderColor: "secondary"
          }  
        }
      },
      animation: {
        typing: "typing 2s steps(20) infinite alternate, blink .7s infinite"
      }
    }
  },
  daisyui: {
      themes: [
        {
          myTheme: theme,
        },
      ],
  },
  plugins: [
    require("daisyui"), 
    require('@tailwindcss/typography'),
  ],
}

