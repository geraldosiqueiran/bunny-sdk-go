import { getBrowser, getPage, disconnectBrowser, outputJSON } from '/home/geraldo/bunny-sdk-go/.claude/skills/chrome-devtools/scripts/lib/browser.js';

async function getOpenAPIData() {
  const browser = await getBrowser();
  const page = await getPage(browser);
  
  await page.goto('https://docs.bunny.net/openapi', { waitUntil: 'networkidle2', timeout: 30000 });
  await new Promise(r => setTimeout(r, 3000));

  const data = await page.evaluate(() => {
    // Extract all text content and links from the page
    const bodyText = document.body?.innerText || '';
    const allLinks = Array.from(document.querySelectorAll('a'));
    const links = allLinks.map(a => ({ text: a.textContent.trim(), href: a.href })).filter(x => x.text.length > 0 && x.text.length < 200);
    
    // Look for download links or spec URLs in the page content
    const urlMatches = bodyText.match(/https?:\/\/[^\s]+\.(json|yaml|yml)/g) || [];
    
    return { bodyText: bodyText.substring(0, 5000), links, urlMatches };
  });

  outputJSON({ success: true, ...data });
  await disconnectBrowser();
}

getOpenAPIData();
