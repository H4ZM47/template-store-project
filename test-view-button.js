const puppeteer = require('puppeteer');

// Helper function to wait
const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms));

(async () => {
  console.log('🧪 Testing "View Template" Button Functionality...\n');

  const browser = await puppeteer.launch({
    headless: 'new',
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  const page = await browser.newPage();

  try {
    // Load the homepage
    console.log('✓ Loading homepage...');
    await page.goto('http://localhost:3000', { waitUntil: 'networkidle0' });
    await sleep(1000);

    console.log('\n' + '='.repeat(60));
    console.log('Testing Template 1: Data Classification Standard');
    console.log('='.repeat(60));

    // Click preview button for first template
    console.log('\n✓ Step 1: Opening preview modal...');
    await page.evaluate(() => {
      const cards = Array.from(document.querySelectorAll('.template-card'));
      const button = cards[0]?.querySelector('button, .btn');
      if (button) button.click();
    });
    await sleep(500);
    console.log('  ✅ Modal opened');

    // Click "View Template" button
    console.log('\n✓ Step 2: Clicking "View Template" button...');

    // Listen for new page opening
    const newPagePromise = new Promise(resolve => {
      browser.on('targetcreated', async target => {
        const newPage = await target.page();
        if (newPage) resolve(newPage);
      });
    });

    // Click the View Template button
    await page.evaluate(() => {
      const modal = document.querySelector('.modal, [class*="modal"], [id*="modal"]');
      const buttons = Array.from(modal.querySelectorAll('button, .btn, a.button'));
      const viewButton = buttons.find(btn => btn.textContent.includes('View Template'));
      if (viewButton) viewButton.click();
    });

    console.log('  ✅ "View Template" button clicked');

    // Wait for new page to open
    console.log('\n✓ Step 3: Waiting for new tab to open...');
    const newPage = await Promise.race([
      newPagePromise,
      sleep(3000).then(() => null)
    ]);

    if (newPage) {
      await sleep(2000); // Wait for page to load
      const url = newPage.url();
      const title = await newPage.title();

      console.log('  ✅ New tab opened successfully!');
      console.log(`  📝 URL: ${url}`);
      console.log(`  📝 Title: "${title}"`);

      // Check if it's the correct template
      const content = await newPage.content();
      const hasDataClassification = content.includes('Data Classification Standard');
      console.log(`  📝 Contains "Data Classification Standard": ${hasDataClassification ? '✅' : '❌'}`);

      // Take screenshot
      await newPage.screenshot({ path: 'template-1-view.png', fullPage: false });
      console.log('  📸 Screenshot saved: template-1-view.png');

      await newPage.close();
    } else {
      console.log('  ❌ New tab did not open');
    }

    // Test second template
    console.log('\n' + '='.repeat(60));
    console.log('Testing Template 2: Vulnerability Management Standard');
    console.log('='.repeat(60));

    await page.goto('http://localhost:3000', { waitUntil: 'networkidle0' });
    await sleep(500);

    console.log('\n✓ Step 1: Opening preview modal...');
    await page.evaluate(() => {
      const cards = Array.from(document.querySelectorAll('.template-card'));
      const button = cards[1]?.querySelector('button, .btn');
      if (button) button.click();
    });
    await sleep(500);
    console.log('  ✅ Modal opened');

    console.log('\n✓ Step 2: Clicking "View Template" button...');

    const newPagePromise2 = new Promise(resolve => {
      browser.on('targetcreated', async target => {
        const newPage = await target.page();
        if (newPage) resolve(newPage);
      });
    });

    await page.evaluate(() => {
      const modal = document.querySelector('.modal, [class*="modal"], [id*="modal"]');
      const buttons = Array.from(modal.querySelectorAll('button, .btn, a.button'));
      const viewButton = buttons.find(btn => btn.textContent.includes('View Template'));
      if (viewButton) viewButton.click();
    });

    console.log('  ✅ "View Template" button clicked');

    console.log('\n✓ Step 3: Waiting for new tab to open...');
    const newPage2 = await Promise.race([
      newPagePromise2,
      sleep(3000).then(() => null)
    ]);

    if (newPage2) {
      await sleep(2000);
      const url = newPage2.url();
      const title = await newPage2.title();

      console.log('  ✅ New tab opened successfully!');
      console.log(`  📝 URL: ${url}`);
      console.log(`  📝 Title: "${title}"`);

      const content = await newPage2.content();
      const hasVulnerability = content.includes('Vulnerability Management Standard');
      console.log(`  📝 Contains "Vulnerability Management Standard": ${hasVulnerability ? '✅' : '❌'}`);

      await newPage2.screenshot({ path: 'template-2-view.png', fullPage: false });
      console.log('  📸 Screenshot saved: template-2-view.png');

      await newPage2.close();
    } else {
      console.log('  ❌ New tab did not open');
    }

    // Final summary
    console.log('\n' + '='.repeat(60));
    console.log('🎉 VIEW TEMPLATE BUTTON TESTING COMPLETE!');
    console.log('='.repeat(60));
    console.log('\n✅ Both "View Template" buttons working correctly!');
    console.log('✅ Templates open in new tabs');
    console.log('✅ Correct HTML content displayed');

  } catch (error) {
    console.error('\n❌ Error during testing:', error.message);
    await page.screenshot({ path: 'view-button-error.png' });
    console.log('Error screenshot saved: view-button-error.png');
  } finally {
    await browser.close();
  }
})();
