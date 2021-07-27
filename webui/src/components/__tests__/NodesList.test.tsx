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

