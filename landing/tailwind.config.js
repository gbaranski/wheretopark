import defaultTheme from 'tailwindcss/defaultTheme';

const colors = {
	red: "#EF2D56",
	green: "#04E762",
	yellow: "#FFC15E",
	orange: "#EC8004",
	blue: "#3FA9F5",
	white: "#FFFFFF",
	darkBlue: "#383F51",
	veryDarkBlue: "#011627",
};

const colorScheme = {
	"primary": colors.blue,
	"secondary": colors.darkBlue,
	"accent": colors.red,
	"neutral": colors.veryDarkBlue,
	"base-100": colors.white,
	"info": colors.blue,
	"success": colors.green,
	"warning": colors.orange,
	"error": colors.red,
};


/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
	theme: {
		extend: {
			colors: {
				primary: colorScheme.primary,
				secondary: colorScheme.secondary,
				accent: colorScheme.accent,
				default: "#101010",
				muted: "rgba(16, 16, 16, 0.66)",
			},
			fontFamily: {
				'sans': ['Inter Variable', ...defaultTheme.fontFamily.sans],
			},
		}
	},
	daisyui: {
		themes: [
			{
				myTheme: colorScheme,
			},
		],
	},
	plugins: [
		require("daisyui"),
		require('@tailwindcss/typography'),
	],
}
