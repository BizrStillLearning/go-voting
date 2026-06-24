import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
    compatibilityDate: '2025-07-15',
    devtools: { enabled: true },

    modules: ["@pinia/nuxt"],

    vite: {
        plugins: [tailwindcss()],
    },

    app: {
        pageTransition: { name: 'page', mode: 'out-in' }
    },

    css: ["~/assets/css/main.css"],

    devServer: {
        host: '0.0.0.0',
        port: 3000
    },

    runtimeConfig: {
        public: {
            apiBase: process.env.API_BASE_URL || 'http://localhost:8080/api'
        }
    }
})