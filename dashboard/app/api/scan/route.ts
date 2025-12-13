import { NextResponse } from 'next/server';
import fs from 'fs';
import path from 'path';

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url);
    const repo = searchParams.get('repo') || '';
    
    // If repo is specified, get from data directory
    if (repo) {
      const isVercel = process.env.VERCEL === '1';
      const baseDir = isVercel ? '/tmp/costguard-data' : path.join(process.cwd(), 'data');
      const repoDir = path.join(baseDir, repo.replace('/', '_'));
      const latestScanFile = path.join(repoDir, 'latest-scan.json');
      
      if (fs.existsSync(latestScanFile)) {
        const scanData = JSON.parse(fs.readFileSync(latestScanFile, 'utf-8'));
        return NextResponse.json(scanData);
      }
      
      return NextResponse.json({
        total_current_cost_usd: 0,
        total_optimal_cost_usd: 0,
        total_potential_savings_usd: 0,
        resources: [],
      });
    }
    
    // Fallback: try local .costguard directory (for local dev)
    const scanPath = path.join(process.cwd(), '../../.costguard/scan.json');
    
    if (fs.existsSync(scanPath)) {
      const scanData = JSON.parse(fs.readFileSync(scanPath, 'utf-8'));
      return NextResponse.json(scanData);
    }
    
    return NextResponse.json({
      total_current_cost_usd: 0,
      total_optimal_cost_usd: 0,
      total_potential_savings_usd: 0,
      resources: [],
    });
  } catch (error) {
    console.error('Error loading scan data:', error);
    return NextResponse.json(
      { error: 'Failed to load scan data' },
      { status: 500 }
    );
  }
}

