# Vercel Dashboard Setup Guide

## Quick Deploy

### 1. Install Vercel CLI

```bash
npm install -g vercel
```

### 2. Login to Vercel

```bash
vercel login
```

### 3. Deploy Dashboard

```bash
cd dashboard
vercel
```

Follow the prompts:
- Set up and deploy? **Yes**
- Which scope? (Select your account)
- Link to existing project? **No**
- Project name? `costguard-dashboard` (or your choice)
- Directory? `./`
- Override settings? **No**

### 4. Set Environment Variables

1. Go to https://vercel.com/dashboard
2. Select your `costguard-dashboard` project
3. Go to **Settings** → **Environment Variables**
4. Add variables:

| Name | Value | Environment |
|------|-------|-------------|
| `KESTRA_API_URL` | `https://your-kestra.com` | Production, Preview, Development |
| `NEXT_PUBLIC_APP_URL` | `https://your-dashboard.vercel.app` | Production |

### 5. Create Webhook Endpoint (Optional)

If you want Kestra to update the dashboard, create:

**File:** `dashboard/app/api/kestra-update/route.ts`

```typescript
import { NextResponse } from 'next/server';

export async function POST(request: Request) {
  try {
    const data = await request.json();
    
    // Store update in database or file system
    // For demo, you could write to .costguard/updates.json
    
    return NextResponse.json({ success: true });
  } catch (error) {
    return NextResponse.json(
      { error: 'Failed to process update' },
      { status: 500 }
    );
  }
}
```

Then set in Kestra secrets:
- `VERCEL_DASHBOARD_WEBHOOK` = `https://your-dashboard.vercel.app/api/kestra-update`

### 6. Redeploy

After setting environment variables:

```bash
vercel --prod
```

Or trigger redeploy from Vercel dashboard.

## Local Development

```bash
cd dashboard
npm install
npm run dev
```

Open http://localhost:3000

## Custom Domain (Optional)

1. Go to **Settings** → **Domains**
2. Add your custom domain
3. Follow DNS configuration instructions

## Monitoring

- **Analytics**: Vercel Analytics (if enabled)
- **Logs**: View in Vercel dashboard → **Deployments** → Click deployment → **Functions** tab
- **Performance**: Vercel Speed Insights

## Troubleshooting

**Build fails:**
- Check Node.js version (should be 18+)
- Verify all dependencies in `package.json`
- Check build logs in Vercel dashboard

**Dashboard shows no data:**
- Ensure `.costguard/scan.json` exists (for local)
- Check API routes are working: `/api/scan`, `/api/decisions`
- Verify environment variables are set

**Webhook not receiving updates:**
- Check webhook URL is correct
- Verify Vercel API token in Kestra
- Check function logs in Vercel dashboard

