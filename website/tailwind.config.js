import defaultTheme from 'tailwindcss/defaultTheme';
import theme from "../svelte/theme";

/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      fontFamily: {
        'sans': ['Poppins', ...defaultTheme.fontFamily.sans],
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

