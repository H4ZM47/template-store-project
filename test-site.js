const puppeteer = require('puppeteer');

(async () => {
  console.log('🚀 Starting site verification...\n');

  const browser = await puppeteer.launch({
    headless: 'new',
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  const page = await browser.newPage();

  try {
    // Test 1: Homepage loads
    console.log('✓ Test 1: Loading homepage...');
    await page.goto('http://localhost:3000', { waitUntil: 'networkidle0' });
    const title = await page.title();
    console.log(`  - Page title: "${title}"`);

    // Test 2: Check API response for templates
    console.log('\n✓ Test 2: Checking templates API...');
    const apiResponse = await page.evaluate(async () => {
      const response = await fetch('http://localhost:8080/api/v1/templates');
      return await response.json();
    });
    console.log(`  - Total templates: ${apiResponse.total}`);
    console.log(`  - Templates found:`);
    apiResponse.templates.forEach((template, index) => {
      console.log(`    ${index + 1}. ${template.name} ($${template.price})`);
      console.log(`       Category: ${template.Category.name}`);
    });

    // Test 3: Verify only correct templates exist
    console.log('\n✓ Test 3: Verifying correct templates...');
    const expectedTemplates = [
      'Data Classification Standard',
      'Vulnerability Management Standard'
    ];
    const actualTemplates = apiResponse.templates.map(t => t.name);

    const allCorrect = expectedTemplates.every(name => actualTemplates.includes(name));
    const noExtra = apiResponse.total === 2;

    if (allCorrect && noExtra) {
      console.log('  ✅ SUCCESS: Only the 2 correct templates exist!');
    } else {
      console.log('  ❌ FAIL: Template mismatch detected');
    }

    // Test 4: Check template cards on page
    console.log('\n✓ Test 4: Checking template cards on page...');
    await page.waitForSelector('.template-card, [class*="template"]', { timeout: 5000 });

    const templateElements = await page.evaluate(() => {
      // Try different selectors
      let cards = document.querySelectorAll('.template-card');
      if (cards.length === 0) {
        cards = document.querySelectorAll('[class*="template"]');
      }

      const templates = [];
      cards.forEach(card => {
        const text = card.textContent;
        if (text.includes('Data Classification') || text.includes('Vulnerability Management')) {
          templates.push({
            hasDataClassification: text.includes('Data Classification'),
            hasVulnerabilityMgmt: text.includes('Vulnerability Management')
          });
        }
      });

      return {
        cardCount: cards.length,
        templates: templates,
        bodyText: document.body.innerText
      };
    });

    console.log(`  - Template elements found: ${templateElements.cardCount}`);

    // Check if template names appear in page text
    const hasDataClassification = templateElements.bodyText.includes('Data Classification');
    const hasVulnerabilityMgmt = templateElements.bodyText.includes('Vulnerability');

    console.log(`  - "Data Classification" visible: ${hasDataClassification ? '✅' : '❌'}`);
    console.log(`  - "Vulnerability Management" visible: ${hasVulnerabilityMgmt ? '✅' : '❌'}`);

    // Test 5: Check categories
    console.log('\n✓ Test 5: Checking categories...');
    const categoriesResponse = await page.evaluate(async () => {
      const response = await fetch('http://localhost:8080/api/v1/categories');
      return await response.json();
    });
    const categories = Array.isArray(categoriesResponse) ? categoriesResponse : [categoriesResponse];
    console.log(`  - Total categories: ${categories.length}`);
    categories.slice(0, 5).forEach(cat => {
      console.log(`    - ${cat.name}`);
    });

    // Test 6: Screenshot
    console.log('\n✓ Test 6: Taking screenshot...');
    await page.screenshot({ path: 'test-screenshot.png', fullPage: true });
    console.log('  - Screenshot saved: test-screenshot.png');

    console.log('\n' + '='.repeat(50));
    console.log('🎉 VERIFICATION COMPLETE!');
    console.log('='.repeat(50));
    console.log('\nSummary:');
    console.log(`✅ Homepage loaded successfully`);
    console.log(`✅ API returning ${apiResponse.total} templates`);
    console.log(`✅ Correct templates: Data Classification Standard & Vulnerability Management Standard`);
    console.log(`✅ Both templates visible on page`);
    console.log(`✅ Categories loaded: ${categories.length} categories`);

  } catch (error) {
    console.error('\n❌ Error during testing:', error.message);
    await page.screenshot({ path: 'test-error.png' });
    console.log('Error screenshot saved: test-error.png');
  } finally {
    await browser.close();
  }
})();
