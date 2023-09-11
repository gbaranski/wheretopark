import defaultTheme from 'tailwindcss/defaultTheme';

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

