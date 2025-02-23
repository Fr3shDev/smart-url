# smart-url

Most URL shorteners simply store a long URL and return a short code that always redirects to the same destination. SmartURL takes it a step further by allowing the creator to specify conditional redirection rules. For example, a marketer might want:

Mobile users (detected via the User-Agent) to be sent to a mobile-optimized page.
Visitors from a specific region (detected via their IP/geolocation) to be sent to a localized landing page.

Everyone else to be redirected to a default URL.
This extra “smart” redirection feature is something you don’t typically see in off-the-shelf URL shorteners.

# Next steps & Enhancements
- Persistent Storage:
Replace the in‑memory map with a database like SQLite or PostgreSQL.

- More Conditions:
Integrate an IP geolocation API (like ip-api.com) to add location-based redirection.

- Analytics:
Log requests and redirection events for later analysis (e.g., how many mobile vs. desktop hits).

- Error Handling & Logging:
Enhance error handling, logging, and add metrics for production-readiness.

- Security Improvements:
Add authentication for URL creation, rate limiting, etc.