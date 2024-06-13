import type { Config } from "tailwindcss";

// default: https://github.com/tailwindlabs/tailwindcss/blob/master/stubs/config.full.js

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      backgroundImage: {
        "gradient-radial": "radial-gradient(var(--tw-gradient-stops))",
        "gradient-conic":
          "conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))",
      },
      colors: {
        bg: {
          DEFAULT: "#0f172a"
        },
        container: {
          DEFAULT: "#1e293b",
          light: "#334155"
        },
        border: {
          DEFAULT: "#334155"
        }
      }
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
};
export default config;
