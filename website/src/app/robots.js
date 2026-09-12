export default function robots() {
  // Must be the site's own origin (matches metadataBase), not the repo URL.
  const baseUrl = "https://chronotrace.org";

  return {
    rules: [
      {
        userAgent: "*",
        allow: "/",
      },
      {
        userAgent: [
          "GPTBot",
          "ChatGPT-User",
          "PerplexityBot",
          "ClaudeBot",
          "Claude-Web",
          "Anthropic-AI",
          "Google-Extended",
          "Applebot-Extended",
          "CCBot",
          "cohere-ai",
          "Bytespider",
        ],
        allow: ["/", "/llms.txt", "/llms-full.txt"],
      },
    ],
    sitemap: `${baseUrl}/sitemap.xml`,
  };
}
