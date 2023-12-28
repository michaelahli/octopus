import { test, expect } from '@playwright/test';

const baseURL = process.env.ENDPOINT || 'http://127.0.0.1:8080';
const uuidRegex = /^[a-f0-9]{8}-[a-f0-9]{4}-[4][a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$/i;


test('if get books then books should be returned', async ({ page }) => {
    const path = '/book';
    const url = `${baseURL}${path}`;

    await page.goto(url);
    await page.waitForLoadState();

    const responseBody = await page.textContent('body');

    let books;
    try {
        books = JSON.parse(responseBody!);
    } catch (error) {
        throw new Error(`Failed to parse JSON response: ${error.message}`);
    }

    expect(Array.isArray(books)).toBeTruthy();
    expect(books.length).toBeGreaterThan(1);

    for (const book of books) {
        expect(book).toHaveProperty('book_id');
        expect(book.book_id).toMatch(uuidRegex);

        expect(book).toHaveProperty('title');
        expect(book.title).not.toBe('');
    }
});
