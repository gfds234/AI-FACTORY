#!/usr/bin/env node
/**
 * Full Validation Script - Zero-Bug Workflow
 * Runs all quality checks before allowing build/deployment
 */

import { execSync } from 'child_process';
import { existsSync } from 'fs';

const colors = {
  reset: '\x1b[0m',
  green: '\x1b[32m',
  red: '\x1b[31m',
  yellow: '\x1b[33m',
  blue: '\x1b[36m'
};

function log(message, color = colors.reset) {
  console.log(`${color}${message}${colors.reset}`);
}

function runCommand(cmd, description) {
  log(`\n🔍 ${description}...`, colors.blue);
  try {
    execSync(cmd, { stdio: 'inherit' });
    log(`✅ ${description} PASSED`, colors.green);
    return true;
  } catch (error) {
    log(`❌ ${description} FAILED`, colors.red);
    return false;
  }
}

function checkFileExists(path, description) {
  log(`\n📄 Checking ${description}...`, colors.blue);
  if (existsSync(path)) {
    log(`✅ ${description} exists`, colors.green);
    return true;
  } else {
    log(`⚠️  ${description} not found (skipping)`, colors.yellow);
    return null; // null = skip, not fail
  }
}

async function main() {
  log('\n' + '='.repeat(60), colors.blue);
  log('  AI FACTORY - ZERO-BUG VALIDATION', colors.blue);
  log('='.repeat(60) + '\n', colors.blue);

  const results = {
    // Required checks
    test: runCommand('npm run test', 'Unit Tests'),
    build: runCommand('npm run build', 'Production Build'),

    // Optional checks (skip if not configured)
    lint: checkFileExists('.eslintrc.json', 'ESLint config')
      ? runCommand('npm run lint', 'Code Linting')
      : null,

    integration: existsSync('tests/integration')
      ? runCommand('npm run test:integration', 'Integration Tests')
      : null,

    e2e: existsSync('e2e') && existsSync('playwright.config.js')
      ? runCommand('npm run test:e2e', 'E2E Tests')
      : null
  };

  // Summary
  log('\n' + '='.repeat(60), colors.blue);
  log('  VALIDATION SUMMARY', colors.blue);
  log('='.repeat(60) + '\n', colors.blue);

  const passed = [];
  const failed = [];
  const skipped = [];

  Object.entries(results).forEach(([name, result]) => {
    if (result === true) passed.push(name);
    else if (result === false) failed.push(name);
    else if (result === null) skipped.push(name);
  });

  if (passed.length > 0) {
    log(`✅ PASSED (${passed.length}): ${passed.join(', ')}`, colors.green);
  }
  if (skipped.length > 0) {
    log(`⚠️  SKIPPED (${skipped.length}): ${skipped.join(', ')}`, colors.yellow);
  }
  if (failed.length > 0) {
    log(`❌ FAILED (${failed.length}): ${failed.join(', ')}`, colors.red);
  }

  log('\n' + '='.repeat(60) + '\n', colors.blue);

  if (failed.length === 0) {
    log('🎉 ALL VALIDATIONS PASSED - Safe to deploy!\n', colors.green);
    process.exit(0);
  } else {
    log('⛔ VALIDATION FAILED - Do NOT deploy!\n', colors.red);
    log('Fix failing tests before proceeding.\n');
    process.exit(1);
  }
}

main().catch(error => {
  log(`\n❌ Validation script error: ${error.message}`, colors.red);
  process.exit(1);
});
