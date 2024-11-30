/** @type {import('tailwindcss').Config} */
module.exports = {
	content: [
		'src/view/**/*.templ',
		'src/view/**/*.go',
		'!**/**docker**/**',
	],
	darkMode: 'class',
	theme: {
		colors: {
			'cblack': '#131720',
			'clight-gray': '#eceff3',
			'clight-light-gray': '#DFE5EB',
			'cwhite': '#CFD0D2',
			'white': '#FFFFFF',
			'cgray': '#42454C',
			'cdark-gray': '#2b2e36'
		},
		container: {
			center: true,
		},
		extend: {
			gridTemplateColumns: {
				'base': '200px 1fr'
			},
			gridTemplateRows: {
				'base': '60px 1fr'
			}
		}
	},
	files : {
		exclude: ["**/**docker**/**","**/.git/**", "**/node_modules/**", "**/.hg/**", "**/.svn/**"]
	},
	plugins: [],
	corePlugins: {
		preflight: true,
	}
}

