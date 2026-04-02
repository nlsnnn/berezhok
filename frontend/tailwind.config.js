/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,jsx}'],
  theme: {
    extend: {
      colors: {
        brand: {
          50:  '#f2f8f4',
          100: '#e0f0e6',
          200: '#c1e1ce',
          300: '#92c9a8',
          400: '#5daa7e',
          500: '#4A7C59',
          600: '#3a6347',
          700: '#2f4f39',
          800: '#27402f',
          900: '#1e3025',
        },
        cream: {
          50:  '#fdfcfa',
          100: '#F5F0E8',
          200: '#ede4d4',
          300: '#ddd0b8',
          400: '#c8b594',
          500: '#b09470',
        },
      },
      fontFamily: {
        sans: ['Manrope', 'system-ui', 'sans-serif'],
      },
      borderRadius: {
        xl:  '1rem',
        '2xl': '1.25rem',
      },
    },
  },
  plugins: [],
}
