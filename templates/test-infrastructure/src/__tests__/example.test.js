/**
 * Example Unit Test
 *
 * This template shows how to write unit tests following TDD principles.
 * Delete this file and create your own tests for your features.
 */

import { describe, it, expect, beforeEach } from 'vitest';

describe('Example Feature', () => {
  // Setup before each test
  beforeEach(() => {
    // Reset state, mock data, etc.
  });

  it('should demonstrate a passing test', () => {
    const result = 2 + 2;
    expect(result).toBe(4);
  });

  it('should test edge cases', () => {
    const emptyArray = [];
    expect(emptyArray.length).toBe(0);
  });

  it('should validate input', () => {
    function validateEmail(email) {
      return email.includes('@');
    }

    expect(validateEmail('test@example.com')).toBe(true);
    expect(validateEmail('invalid')).toBe(false);
  });
});

/**
 * TDD Workflow:
 *
 * 1. Write test that fails:
 *    it('should do the thing', () => {
 *      const result = doTheThing();
 *      expect(result).toBe('expected');
 *    });
 *
 * 2. Run test - see it FAIL ❌
 *    npm run test:watch
 *
 * 3. Write minimum code to pass:
 *    function doTheThing() {
 *      return 'expected';
 *    }
 *
 * 4. Run test - see it PASS ✅
 *
 * 5. Refactor (keep tests green)
 */
