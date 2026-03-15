/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      fontFamily: {
        display: ['"Syne"', 'sans-serif'],
        body: ['"DM Sans"', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'monospace'],
      },
      colors: {
        brand: {
          DEFAULT: '#0EA5A0',
          light: '#14C4BF',
          dark: '#0B8580',
        },
        jade: '#10B981',
        danger: '#EF4444',
        surface: 'var(--surface)',
        bg: 'var(--bg)',
      },
      animation: {
        'fade-in': 'fadeIn 0.35s ease both',
        'slide-up': 'slideUp 0.4s cubic-bezier(0.16, 1, 0.3, 1) both',
        'pulse-brand': 'pulse-brand 2s ease-in-out infinite',
      },
    },
  },
  plugins: [],
}
