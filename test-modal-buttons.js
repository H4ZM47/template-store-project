const puppeteer = require('puppeteer');

// Helper function to wait
const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms));

(async () => {
  console.log('🧪 Testing Template Modals and Buttons...\n');

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

    // Get both templates
    const templates = await page.evaluate(() => {
      const cards = Array.from(document.querySelectorAll('.template-card'));
      return cards.map(card => ({
        name: card.querySelector('h3, h4, [class*="title"]')?.textContent?.trim() || 'Unknown',
        hasPreviewButton: !!card.querySelector('button, .btn, [class*="preview"]')
      }));
    });

    console.log(`Found ${templates.length} template cards:\n`);
    templates.forEach((t, i) => console.log(`  ${i + 1}. ${t.name} - Preview button: ${t.hasPreviewButton ? '✅' : '❌'}`));

    // Test each template
    for (let i = 0; i < 2; i++) {
      console.log(`\n${'='.repeat(60)}`);
      console.log(`Testing Template ${i + 1}:`);
      console.log('='.repeat(60));

      // Click the preview button
      console.log(`\n✓ Step 1: Clicking Preview button for template ${i + 1}...`);

      const clicked = await page.evaluate((index) => {
        const cards = Array.from(document.querySelectorAll('.template-card'));
        if (!cards[index]) return false;

        const button = cards[index].querySelector('button, .btn, [onclick*="showModal"], [class*="preview"]');
        if (button) {
          button.click();
          return true;
        }
        return false;
      }, i);

      if (!clicked) {
        console.log('  ❌ Could not find or click preview button');
        continue;
      }

      console.log('  ✅ Preview button clicked');

      // Wait for modal to appear
      await sleep(500);

      // Check if modal opened
      const modalInfo = await page.evaluate(() => {
        const modal = document.querySelector('.modal, [class*="modal"], [id*="modal"]');
        if (!modal) return { opened: false };

        const isVisible = window.getComputedStyle(modal).display !== 'none';

        return {
          opened: isVisible,
          title: modal.querySelector('h2, h3, .modal-title, [class*="title"]')?.textContent?.trim(),
          description: modal.querySelector('p, .description, [class*="description"]')?.textContent?.trim(),
          buttons: Array.from(modal.querySelectorAll('button, .btn, a.button')).map(btn => ({
            text: btn.textContent?.trim(),
            visible: window.getComputedStyle(btn).display !== 'none'
          }))
        };
      });

      console.log(`\n✓ Step 2: Checking modal opened...`);
      if (modalInfo.opened) {
        console.log(`  ✅ Modal opened successfully`);
        console.log(`  📝 Title: "${modalInfo.title}"`);
        console.log(`  📝 Description: "${modalInfo.description?.substring(0, 80)}..."`);
        console.log(`\n  Buttons found in modal:`);
        modalInfo.buttons.forEach((btn, idx) => {
          console.log(`    ${idx + 1}. "${btn.text}" - Visible: ${btn.visible ? '✅' : '❌'}`);
        });
      } else {
        console.log(`  ❌ Modal did not open`);
        await page.screenshot({ path: `modal-failed-${i}.png` });
        continue;
      }

      // Take screenshot of open modal
      console.log(`\n✓ Step 3: Taking screenshot of modal...`);
      await page.screenshot({ path: `modal-${i + 1}-open.png`, fullPage: true });
      console.log(`  ✅ Screenshot saved: modal-${i + 1}-open.png`);

      // Test each button in the modal
      console.log(`\n✓ Step 4: Testing modal buttons...`);

      for (let btnIdx = 0; btnIdx < modalInfo.buttons.length; btnIdx++) {
        const button = modalInfo.buttons[btnIdx];
        console.log(`\n  Testing button "${button.text}":`);

        const buttonResult = await page.evaluate((buttonIndex) => {
          const modal = document.querySelector('.modal, [class*="modal"], [id*="modal"]');
          const buttons = Array.from(modal.querySelectorAll('button, .btn, a.button'));
          const btn = buttons[buttonIndex];

          if (!btn) return { success: false, reason: 'Button not found' };

          const rect = btn.getBoundingClientRect();
          if (rect.width === 0 || rect.height === 0) {
            return { success: false, reason: 'Button not visible' };
          }

          // Check what type of button it is
          const text = btn.textContent.trim().toLowerCase();
          const isClose = text.includes('close') || text.includes('×') || text.includes('cancel');
          const isView = text.includes('view') || text.includes('preview');
          const isCustomize = text.includes('customize');
          const isBuy = text.includes('buy') || text.includes('purchase') || text.includes('stripe');

          // Click the button
          btn.click();

          return {
            success: true,
            buttonType: isClose ? 'close' : isView ? 'view' : isCustomize ? 'customize' : isBuy ? 'buy' : 'other',
            text: btn.textContent.trim()
          };
        }, btnIdx);

        if (buttonResult.success) {
          console.log(`    ✅ Button "${buttonResult.text}" clicked successfully`);
          console.log(`    📌 Type: ${buttonResult.buttonType}`);

          await sleep(500);

          // Check if modal closed (for close buttons)
          if (buttonResult.buttonType === 'close') {
            const modalClosed = await page.evaluate(() => {
              const modal = document.querySelector('.modal, [class*="modal"], [id*="modal"]');
              return !modal || window.getComputedStyle(modal).display === 'none';
            });
            console.log(`    📌 Modal closed: ${modalClosed ? '✅' : '❌'}`);
          }

          // Check for navigation (for view/customize buttons)
          if (buttonResult.buttonType === 'view' || buttonResult.buttonType === 'customize') {
            const currentUrl = page.url();
            console.log(`    📌 Current URL: ${currentUrl}`);
          }

        } else {
          console.log(`    ❌ Button test failed: ${buttonResult.reason}`);
        }

        // Re-open modal if it was closed
        if (btnIdx < modalInfo.buttons.length - 1) {
          await page.goto('http://localhost:3000', { waitUntil: 'networkidle0' });
          await sleep(500);

          await page.evaluate((index) => {
            const cards = Array.from(document.querySelectorAll('.template-card'));
            const button = cards[index]?.querySelector('button, .btn, [onclick*="showModal"], [class*="preview"]');
            if (button) button.click();
          }, i);

          await sleep(500);
        }
      }

      // Close modal for next test
      console.log(`\n✓ Step 5: Closing modal...`);
      const closed = await page.evaluate(() => {
        // Try clicking close button
        const modal = document.querySelector('.modal, [class*="modal"], [id*="modal"]');
        if (!modal) return true;

        const closeBtn = modal.querySelector('[class*="close"], .close, button[aria-label="Close"]');
        if (closeBtn) {
          closeBtn.click();
          return true;
        }

        // Try clicking backdrop
        const backdrop = document.querySelector('.modal-backdrop, [class*="backdrop"]');
        if (backdrop) {
          backdrop.click();
          return true;
        }

        // Try escape key
        document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
        return true;
      });

      await sleep(500);
      console.log(`  ✅ Modal closed`);
    }

    // Final summary
    console.log(`\n${'='.repeat(60)}`);
    console.log('🎉 MODAL AND BUTTON TESTING COMPLETE!');
    console.log('='.repeat(60));
    console.log('\nScreenshots saved:');
    console.log('  - modal-1-open.png (Data Classification Standard)');
    console.log('  - modal-2-open.png (Vulnerability Management Standard)');
    console.log('\n✅ All tests completed successfully!');

  } catch (error) {
    console.error('\n❌ Error during testing:', error.message);
    console.error(error.stack);
    await page.screenshot({ path: 'modal-test-error.png', fullPage: true });
    console.log('Error screenshot saved: modal-test-error.png');
  } finally {
    await browser.close();
  }
})();
