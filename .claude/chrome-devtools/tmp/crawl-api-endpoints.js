import { getBrowser, getPage, disconnectBrowser, outputJSON } from '/home/geraldo/bunny-sdk-go/.claude/skills/chrome-devtools/scripts/lib/browser.js';

async function crawlEndpoints() {
  const browser = await getBrowser();
  const page = await getPage(browser);

  // Start from each API section's first endpoint page and follow Next links
  const startPages = [
    // Core API - start from first endpoint in each group
    'https://docs.bunny.net/api-reference/core/pull-zones/list-pull-zones',
    'https://docs.bunny.net/api-reference/core/storage-zones/list-storage-zones',
    'https://docs.bunny.net/api-reference/core/dns/list-dns-zones',
    'https://docs.bunny.net/api-reference/core/stream-video-library/list-video-libraries',
    'https://docs.bunny.net/api-reference/core/statistics/get-statistics',
    'https://docs.bunny.net/api-reference/core/purge/purge-url',
    'https://docs.bunny.net/api-reference/core/api-keys/list-api-keys',
    'https://docs.bunny.net/api-reference/core/countries/country-list',
    'https://docs.bunny.net/api-reference/core/region/region-list',
    'https://docs.bunny.net/api-reference/core/user-operations/close-account',
    'https://docs.bunny.net/api-reference/core/search/global-search-endpoint',
    // Storage file operations
    'https://docs.bunny.net/api-reference/storage/browse-files/list-files',
    'https://docs.bunny.net/api-reference/storage/manage-files/download-file',
    // Stream
    'https://docs.bunny.net/api-reference/stream/manage-collections/get-collection',
    'https://docs.bunny.net/api-reference/stream/manage-videos/list-videos',
    // Shield
    'https://docs.bunny.net/api-reference/shield/waf/get-shieldwafcustom-rules',
    'https://docs.bunny.net/api-reference/shield/accesslists/get-shieldshield-zone-access-lists',
    // Scripting
    'https://docs.bunny.net/api-reference/scripting/edge-script/list-edge-scripts',
    // Magic Containers
    'https://docs.bunny.net/api-reference/magic-containers/applications/list-applications',
  ];

  const visited = new Set();
  const endpoints = [];

  for (const startUrl of startPages) {
    let currentUrl = startUrl;
    let maxPages = 50; // Safety limit per section
    
    while (currentUrl && maxPages > 0 && !visited.has(currentUrl)) {
      visited.add(currentUrl);
      maxPages--;

      try {
        await page.goto(currentUrl, { waitUntil: 'networkidle2', timeout: 20000 });
        await new Promise(r => setTimeout(r, 1500));

        const data = await page.evaluate(() => {
          const title = document.title || '';
          const bodyText = document.body?.innerText || '';
          
          // Look for HTTP method badge and path in the page content
          // Typically shown as "GET /path/to/endpoint" or similar
          const methodMatch = bodyText.match(/(GET|POST|PUT|PATCH|DELETE|HEAD)\s+(\/[^\s\n]+)/);
          const method = methodMatch ? methodMatch[1] : '';
          const path = methodMatch ? methodMatch[2] : '';
          
          // Find description - usually the first paragraph after the title
          const h1 = document.querySelector('h1');
          const desc = h1 ? h1.textContent.trim() : title.split(' - ')[0];

          // Find Next link
          const allLinks = Array.from(document.querySelectorAll('a'));
          const nextLink = allLinks.find(a => {
            const text = a.textContent.trim();
            return text.endsWith('Next') && a.href.includes('/api-reference/');
          });

          return {
            title: desc,
            method,
            path,
            nextUrl: nextLink ? nextLink.href : null,
            url: window.location.href
          };
        });

        if (data.method && data.path) {
          endpoints.push({
            method: data.method,
            path: data.path,
            title: data.title,
            url: data.url
          });
        }

        // Follow Next link if it's in a different section or hasn't been visited
        currentUrl = data.nextUrl;
        if (currentUrl && visited.has(currentUrl)) {
          currentUrl = null;
        }
      } catch (e) {
        currentUrl = null;
      }
    }
  }

  outputJSON({ success: true, totalEndpoints: endpoints.length, endpoints });
  await disconnectBrowser();
}

crawlEndpoints();
