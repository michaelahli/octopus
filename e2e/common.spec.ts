import { test, expect } from '@playwright/test';

const baseURL = process.env.ENDPOINT || 'http://127.0.0.1:8080';
const runtimeEnvironment = process.env.RUNTIME_ENV || 'staging';

test('if request /:path then /:path should be reflected', async ({ page }) => {
    const path = '/:path';
    const url = `${baseURL}${path}`;
    await page.goto(url);
    await page.waitForLoadState();
    const bodyText = await page.textContent('body');
    await expect(bodyText).toContain(`Hello from ${runtimeEnvironment} ${path}`);
});
