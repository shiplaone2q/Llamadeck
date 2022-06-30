import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { nodeHealthService } from '@/lib/nodeHealthService';

// real service + real apiCall, only fetch is mocked
describe('nodeHealthService', () => {
  const okBody = (status: string) =>
    new Response(`{"status":"${status}","checked_at":1}\n`, {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    });

  beforeEach(() => {
    global.fetch = vi.fn(() =>
      Promise.resolve(okBody('healthy')),
    ) as unknown as typeof fetch;
  });

  afterEach(() => {
    nodeHealthService.destroy();
  });

  it('maps a healthy response to a healthy state', async () => {
    const health = await nodeHealthService.performHealthCheck('main');
    expect(health.state).toBe('healthy');
    expect(health.error).toBeUndefined();
  });

  it('maps an unreachable response to an unreachable state with error', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve(
        new Response(
          '{"status":"unreachable","error":"connection refused","checked_at":1}\n',
          { status: 200, headers: { 'Content-Type': 'application/json' } },
        ),
      ),
    ) as unknown as typeof fetch;
