import { NextResponse } from 'next/server';
import fs from 'fs';
import path from 'path';

// Ensure this route is publicly accessible
export const runtime = 'nodejs';
export const dynamic = 'force-dynamic';

interface DashboardUpdate {
  repo_owner: string;
  repo_name: string;
  repo_full_name: string;
  scan_data: any;
  decision_data: any;
  pr_url?: string;
  pr_number?: number;
  timestamp: string;
}

export async function POST(request: Request) {
  try {
    // Add CORS headers
    const headers = {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'POST, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization',
    };

    // Handle OPTIONS request for CORS
    if (request.method === 'OPTIONS') {
      return new NextResponse(null, { status: 200, headers });
    }

    const data: DashboardUpdate = await request.json();
    
    // Use /tmp directory for Vercel (writable in serverless functions)
    // In local dev, use data directory
    const isVercel = process.env.VERCEL === '1';
    const baseDir = isVercel ? '/tmp/costguard-data' : path.join(process.cwd(), 'data');
    
    // Create data directory
    if (!fs.existsSync(baseDir)) {
      fs.mkdirSync(baseDir, { recursive: true });
    }
    
    // Create repo-specific directory
    const repoDir = path.join(baseDir, data.repo_full_name.replace('/', '_'));
    if (!fs.existsSync(repoDir)) {
      fs.mkdirSync(repoDir, { recursive: true });
    }
    
    // Save scan data
    const scanFile = path.join(repoDir, `scan-${Date.now()}.json`);
    fs.writeFileSync(scanFile, JSON.stringify({
      ...data.scan_data,
      repo_full_name: data.repo_full_name,
      timestamp: data.timestamp,
    }, null, 2));
    
    // Save decision data
    const decisionFile = path.join(repoDir, `decision-${Date.now()}.json`);
    fs.writeFileSync(decisionFile, JSON.stringify({
      ...data.decision_data,
      repo_full_name: data.repo_full_name,
      pr_url: data.pr_url || '',
      pr_number: data.pr_number || 0,
      timestamp: data.timestamp,
    }, null, 2));
    
    // Update latest scan for this repo
    const latestScanFile = path.join(repoDir, 'latest-scan.json');
    fs.writeFileSync(latestScanFile, JSON.stringify({
      ...data.scan_data,
      repo_full_name: data.repo_full_name,
      timestamp: data.timestamp,
    }, null, 2));
    
    // Update latest decision for this repo
    const latestDecisionFile = path.join(repoDir, 'latest-decision.json');
    fs.writeFileSync(latestDecisionFile, JSON.stringify({
      ...data.decision_data,
      repo_full_name: data.repo_full_name,
      pr_url: data.pr_url || '',
      pr_number: data.pr_number || 0,
      timestamp: data.timestamp,
    }, null, 2));
    
    return NextResponse.json({ 
      success: true,
      repo: data.repo_full_name,
      timestamp: data.timestamp,
    }, { headers });
  } catch (error: any) {
    console.error('Error processing dashboard update:', error);
    return NextResponse.json(
      { 
        error: 'Failed to process update',
        message: error.message,
        details: process.env.VERCEL ? 'Running on Vercel' : 'Running locally'
      },
      { status: 500}
    );
  }
}

