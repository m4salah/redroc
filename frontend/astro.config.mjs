import { defineConfig } from 'astro/config';
import nodejs from "@astrojs/node";
import tailwind from "@astrojs/tailwind";

import react from "@astrojs/react";

// https://astro.build/config
export default defineConfig({
  adapter: nodejs({
    mode: "standalone"
  }),
  output: "server",
  server: {
    port: 3000
  },
  integrations: [tailwind({
    applyBaseStyles: false,
  }), react()]
});