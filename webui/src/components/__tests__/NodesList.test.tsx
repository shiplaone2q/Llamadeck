import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import NodesList from '@/components/nodes/NodesList';
import { AuthProvider } from '@/contexts/AuthContext';
import { nodesApi } from '@/lib/api';
import { nodeHealthService } from '@/lib/nodeHealthService';
import type { NodeHealth } from '@/types/node';

vi.mock('@/lib/api', () => ({
  nodesApi: {
    list: vi.fn(),
    get: vi.fn(),
    getHealth: vi.fn(),
  },
  authApi: {
    whoami: vi.fn(() =>
      Promise.resolve({
        authenticated: false,
        oidc_enabled: false,
        user: null,
      }),
    ),
  },
}));

vi.mock('@/lib/nodeHealthService', () => ({
  nodeHealthService: {
    subscribe: vi.fn(() => () => {}),
  },
}));

vi.mock('@/contexts/ConfigContext', () => ({
  useConfig: () => ({
    config: { local_node: 'main' },
    isLoading: false,
    error: null,
  }),
}));

const mockNodes = {
  'worker-2': { address: 'http://worker-2:8080' },
  main: { address: '' },
  'worker-1': { address: 'http://worker-1:8080' },
};

const healthy: NodeHealth = {
  state: 'healthy',
  latencyMs: 12,
  lastChecked: new Date(),
};

function renderNodesList() {
  return render(
    <AuthProvider>
      <NodesList />
    </AuthProvider>,
  );
}

describe('NodesList', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    window.sessionStorage.setItem(
      'llamadeck_management_key',
      'test-api-key-123',
    );
    global.fetch = vi.fn(() =>
      Promise.resolve(new Response(null, { status: 200 })),
    );
    // Radix tooltips need these in jsdom
    window.HTMLElement.prototype.hasPointerCapture = vi.fn(() => false);
    window.HTMLElement.prototype.setPointerCapture = vi.fn();
    window.HTMLElement.prototype.releasePointerCapture = vi.fn();
    class ResizeObserverMock {
      observe(): void {}
      unobserve(): void {}
      disconnect(): void {}
    }
    window.ResizeObserver =
      ResizeObserverMock as unknown as typeof ResizeObserver;

    vi.mocked(nodesApi.list).mockResolvedValue(mockNodes);
    vi.mocked(nodeHealthService.subscribe).mockImplementation(
      (_name: string, cb: (health: NodeHealth) => void) => {
        cb(healthy);
        return () => {};
      },
    );
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('shows loading spinner while nodes are being fetched', () => {
    vi.mocked(nodesApi.list).mockImplementation(
      () => new Promise((resolve) => setTimeout(() => resolve(mockNodes), 100)),
    );

    renderNodesList();

    expect(screen.getByText('Loading nodes...')).toBeInTheDocument();
  });

  it('displays error message when node loading fails', async () => {
    vi.mocked(nodesApi.list).mockRejectedValue(
      new Error('Failed to connect to server'),
