/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        dark: {
          bg: '#121212',
          surface: '#1e1e1e',
          border: '#333333',
          text: '#ffffff',
          muted: '#a0a0a0'
        },
        primary: {
          DEFAULT: '#3b82f6',
          hover: '#2563eb'
        }
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms')
  ],
}
