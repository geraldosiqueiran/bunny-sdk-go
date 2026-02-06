import { getBrowser, getPage, disconnectBrowser, outputJSON } from '/home/geraldo/bunny-sdk-go/.claude/skills/chrome-devtools/scripts/lib/browser.js';

async function extractEndpoints() {
  const browser = await getBrowser();
  const page = await getPage(browser);

  // First, navigate to the main API reference to get all sub-page links
  const sections = [
    'https://docs.bunny.net/api-reference/core',
    'https://docs.bunny.net/api-reference/storage',
    'https://docs.bunny.net/api-reference/stream',
    'https://docs.bunny.net/api-reference/shield',
    'https://docs.bunny.net/api-reference/scripting',
    'https://docs.bunny.net/api-reference/magic-containers/overview',
  ];

  const allEndpointLinks = [];

  for (const sectionUrl of sections) {
    await page.goto(sectionUrl, { waitUntil: 'networkidle2', timeout: 30000 });
    await new Promise(r => setTimeout(r, 2000));

    const links = await page.evaluate(() => {
      const allLinks = Array.from(document.querySelectorAll('a'));
      return allLinks
        .filter(a => a.href && a.href.includes('/api-reference/'))
        .map(a => ({
          text: a.textContent.trim(),
          href: a.href
        }))
        .filter(x => x.text.length > 0 && x.text.length < 100);
    });

    allEndpointLinks.push({ section: sectionUrl, links });
  }

  // Deduplicate links
  const uniqueLinks = new Map();
  for (const section of allEndpointLinks) {
    for (const link of section.links) {
      if (!uniqueLinks.has(link.href)) {
        uniqueLinks.set(link.href, link.text);
      }
    }
  }

  const sortedLinks = Array.from(uniqueLinks.entries())
    .map(([href, text]) => ({ href, text }))
    .sort((a, b) => a.href.localeCompare(b.href));

  outputJSON({ success: true, totalLinks: sortedLinks.length, links: sortedLinks });
  await disconnectBrowser();
}

extractEndpoints();
