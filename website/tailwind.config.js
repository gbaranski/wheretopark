import defaultTheme from 'tailwindcss/defaultTheme';
import theme from "../svelte/theme";

/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      fontFamily: {
        'sans': ['Poppins', ...defaultTheme.fontFamily.sans],
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

