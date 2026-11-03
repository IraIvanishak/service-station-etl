// Base URL of the Go API. Set VITE_API_URL in .env.local to point elsewhere.
export const API_URL = (import.meta.env.VITE_API_URL || 'http://localhost:3000').replace(/\/$/, '')
