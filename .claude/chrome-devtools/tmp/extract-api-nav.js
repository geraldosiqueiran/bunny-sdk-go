import { getBrowser, getPage, disconnectBrowser, outputJSON } from '/home/geraldo/bunny-sdk-go/.claude/skills/chrome-devtools/scripts/lib/browser.js';

async function extractAPISidebar() {
  const browser = await getBrowser();
  const page = await getPage(browser);

  await page.goto('https://docs.bunny.net/api-reference/stream', { waitUntil: 'networkidle2', timeout: 30000 });
  await new Promise(r => setTimeout(r, 3000));

  const data = await page.evaluate(() => {
    const allLinks = Array.from(document.querySelectorAll('a'));
    const navLinks = allLinks
      .filter(a => a.href && (a.href.includes('/api-reference/') || a.href.includes('/reference/')))
      .map(a => ({
        text: a.textContent.trim(),
        href: a.href
      }))
      .filter(x => x.text.length > 0);
    
    // Get page section headers and endpoint-like content
    const mainText = document.body?.innerText || '';
    
    return { navLinks, url: window.location.href, title: document.title, bodyLength: mainText.length };
  });

  outputJSON({ success: true, ...data });
  await disconnectBrowser();
}

extractAPISidebar();
