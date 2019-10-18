import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from '@/App'
import { InstancesProvider } from '@/contexts/InstancesContext'
import { instancesApi } from '@/lib/api'
import {BackendType, type Instance } from '@/types/instance'
import { AuthProvider } from '@/contexts/AuthContext'

// Mock the API
vi.mock('@/lib/api', () => ({
  instancesApi: {
    list: vi.fn(),
    get: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    start: vi.fn(),
    stop: vi.fn(),
    restart: vi.fn(),
    delete: vi.fn(),
    getHealth: vi.fn(),
  },
  nodesApi: {
    list: vi.fn(() => Promise.resolve({})),
  },
  serverApi: {
    getHelp: vi.fn(),
    getVersion: vi.fn(),
    getDevices: vi.fn(),
  }
}))

// Mock health service to avoid real network calls
vi.mock('@/lib/healthService', () => ({
  healthService: {
    subscribe: vi.fn(() => () => {}),
    refreshHealth: vi.fn(() => Promise.resolve()),
    checkHealthAfterOperation: vi.fn(),
    performHealthCheck: vi.fn(() => Promise.resolve({
      state: 'ready',
      instanceStatus: 'running',
      lastChecked: new Date(),
      source: 'http'
    })),
  },
  checkHealth: vi.fn(() => Promise.resolve({
    state: 'ready',
    instanceStatus: 'running',
    lastChecked: new Date(),
    source: 'http'
  })),
}))

// Mock the ConfigContext helper hooks
vi.mock('@/hooks/useConfig', () => ({
  useInstanceDefaults: () => ({
    idleTimeout: 30,
    autoRestart: true,
    maxRestarts: 3,
    restartDelay: 5,
    onDemandStart: false,
  }),
  useBackendSettings: () => ({
    command: '/usr/bin/llama-server',
    dockerEnabled: false,
    dockerImage: '',
  }),
  useCustomBackends: () => ({}),
  useCustomBackendSettings: () => null,
}))

function renderApp() {
  return render(
    <AuthProvider>
      <InstancesProvider>
        <App />
