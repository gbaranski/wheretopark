import type { BaseTranslation } from "../i18n-types"

const en = {
	link: "{0}",
	heading: "Your parking assistant",
	subheading: "with the help of <b>AI</b> we'll find you an <b>available spot</b> nearby",
	or: "or",
	openAppInBrowser: "open app in the browser",
	minutes: "{0} minutes",
	tonsOf: "{1} tons of {0}",
	problem: {
		what: "What problem are we solving?",
		averageTimeWasted: "Average <b>time wasted</b> each time you want to park.",
		ofDriversCirculate: "Of drivers <b>circulate around</b> parking lots when seeking a spot.",
		airPollutionAnnualContribution: "Annual contribution to the <b>air pollution</b> in {0} alone",
	},
	OUR_SOLUTION: "Our solution",
	OUR_SOLUTION_TEXT: "Easy access to information about available parking lots nearby",
	POWERED_BY_AI: "Powered by AI",
	OUR_SYSTEM_MARKS: `
				Our system automatically marks
				<b>
					<span class="font-bold text-green-600">green</span>
					vacant 
				</b>
				spaces, and
				<b>
					<span class="text-red-600">red</span>
					occupied.
				</b>
				`,
	WE_HAVE_PARKNG_LOTS: "We have parking lots in",
	AWARDS: "Awards",
	TALKS_ABOUT_US: "Talks about us",
	BENEFITS:"Benefits",

	FOR_DRIVERS: "For drivers",
	FOR_DRIVERS_SAVINGS: "<b>Money savings</b> of around {0} per each parking session (based on costs of running the car and wasted time)",
	FOR_DRIVERS_DISABILITIES: "<b>People with disabilities</b> can easily check which parking has specially designated spaces for them, and whether these spaces are free.",
	FOR_DRIVERS_PREDICTIONS: "<b>Predictions of availability</b>. The user has access to the statistics of the occupancy of places on parking lots at a given time and on a given day. For example on Friday at 5 pm, the parking lot is usually 80% full, and on Saturday at 2 pm, it is 50% full.",

	FOR_OPERATORS: "For parking operators",
	FOR_OPERATORS_PROMOTION: "<b>In-app promotion</b> - more customers.",
	FOR_OPERATORS_RESERVATIONS: "<b>Reserving parking spaces</b> directly in the app, without need of phone calling.",

	FOR_CITIES: "For cities",
	FOR_CITIES_ANALYSIS: "<b>Trend analysis</b> - will help cities choosing the right place for a new parking lot, or managing the prices.",
	FOR_CITIES_INCORRECTLY_PARKED: "Our system can <b>detect illegally parked vehicle</b> and notify the authorities.",
	FOR_CITIES_DECREASED_TRAFFIC: "<b>Decreased car traffic</b> - in some cities up to {0} of traffic is caused by drivers seeking a parking spot.",
	FOR_CITIES_OPTIMISED_DISTRIBUTION: "<b>Optimised distribution of cars</b> - drivers could start using less known, or recently opened parking lots, thus reducing the load on the most demanded ones.",
	
	ARE_YOU_PARKING_OPERATOR: "Are you a parking operator?",
	CONTACT_US: "Contact us",
	
	GENERAL_CONTACT: "General Contact",
	OUR_TEAM: "Our Team",
	FAQ: "Frequently Asked Questions",
	FAQ_APP_AVAILABILITY_Q: "When the app will be available?",
	FAQ_APP_AVAILABILITY_A_IOS: "<b>iOS</b> app is available on {0} since 21st of November 2022.",
	FAQ_APP_AVAILABILITY_A_ANDROID: "<b>Android</b> app is available on {0}",
	FAQ_APP_AVAILABILITY_A_WEB: "<b>Web</b> app is available at {0}",
	
	FAQ_WHAT_PARKING_LOTS_Q: "What parking lots do you support?",
	

	
} satisfies BaseTranslation

export default en
