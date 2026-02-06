import { getBrowser, getPage, disconnectBrowser, outputJSON } from '/home/geraldo/bunny-sdk-go/.claude/skills/chrome-devtools/scripts/lib/browser.js';

async function getOpenAPILinks() {
  const browser = await getBrowser();
  const page = await getPage(browser);
  
  await page.goto('https://docs.bunny.net/openapi', { waitUntil: 'networkidle2', timeout: 30000 });
  await new Promise(r => setTimeout(r, 2000));

  const data = await page.evaluate(() => {
    const allLinks = Array.from(document.querySelectorAll('a'));
    return allLinks
      .filter(a => a.href && (a.href.includes('.json') || a.href.includes('.yaml') || a.href.includes('openapi') || a.href.includes('swagger')))
      .map(a => ({ text: a.textContent.trim(), href: a.href }));
  });

  outputJSON({ success: true, links: data });
  await disconnectBrowser();
}

getOpenAPILinks();
