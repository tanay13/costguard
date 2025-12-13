# Vercel Data Storage Solution

## Problem

Vercel serverless functions have a read-only filesystem (except `/tmp`). The `/tmp` directory is ephemeral and data is lost between function invocations.

## Current Solution (Temporary)

The dashboard currently uses `/tmp/costguard-data` for storage on Vercel. This works but data is **not persistent** - it's lost when the function cold starts.

## Recommended Solutions

### Option 1: Vercel KV (Redis) - Recommended

Vercel KV is a Redis-compatible database perfect for this use case.

**Setup:**
1. Install Vercel KV:
   ```bash
   npm install @vercel/kv
   ```

2. Create KV database in Vercel dashboard

3. Update API routes to use KV instead of filesystem

**Example:**
```typescript
import { kv } from '@vercel/kv';

// Store data
await kv.set(`repo:${repoFullName}:latest-scan`, scanData);
await kv.set(`repo:${repoFullName}:latest-decision`, decisionData);

// Retrieve data
const scanData = await kv.get(`repo:${repoFullName}:latest-scan`);
```

### Option 2: Vercel Postgres

For more complex queries and relationships.

### Option 3: External Database

Use MongoDB, PostgreSQL, or other database hosted elsewhere.

### Option 4: GitHub as Storage

Store data in a GitHub repository (gist or repo).

## Quick Fix for Now

The current implementation uses `/tmp` which works but data is ephemeral. For a hackathon demo, this is acceptable if you:
- Accept that data resets on cold starts
- Or redeploy frequently to keep functions warm
- Or use it primarily for local development

## Migration Guide

To migrate to Vercel KV:

1. **Install KV:**
   ```bash
   cd dashboard
   npm install @vercel/kv
   ```

2. **Create KV database:**
   - Go to Vercel dashboard
   - Create KV database
   - Copy connection details

3. **Update environment variables:**
   - Add `KV_REST_API_URL`
   - Add `KV_REST_API_TOKEN`

4. **Update API routes:**
   - Replace filesystem operations with KV operations
   - See example above

