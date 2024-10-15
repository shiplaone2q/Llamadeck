import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import InstanceList from '@/components/InstanceList'
import { InstancesProvider } from '@/contexts/InstancesContext'
import { instancesApi } from '@/lib/api'
import { BackendType, type Instance } from '@/types/instance'
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
  }
}))

// Mock health service
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

function renderInstanceList(editInstance = vi.fn()) {
  return render(
    <AuthProvider>
      <InstancesProvider>
        <InstanceList editInstance={editInstance} />
      </InstancesProvider>
    </AuthProvider>
  )
}

describe('InstanceList - State Management and UI Logic', () => {

  const mockEditInstance = vi.fn()

  const mockInstances: Instance[] = [
    { id: 1, name: 'instance-1', status: 'stopped', options: { backend_type: BackendType.LLAMA_CPP, backend_options: { model: 'model1.gguf' } } },
    { id: 1, name: 'instance-2', status: 'running', options: { backend_type: BackendType.LLAMA_CPP, backend_options: { model: 'model2.gguf' } } },
    { id: 1, name: 'instance-3', status: 'stopped', options: { backend_type: BackendType.LLAMA_CPP, backend_options: { model: 'model3.gguf' } } }
  ]

  const DUMMY_API_KEY = 'test-api-key-123'

  beforeEach(() => {
    vi.clearAllMocks()
    window.sessionStorage.setItem('llamadeck_management_key', DUMMY_API_KEY)
    global.fetch = vi.fn(() => Promise.resolve(new Response(null, { status: 200 })))
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('Loading State', () => {
    it('shows loading spinner while instances are being fetched', () => {
      // Mock a delayed response to test loading state
      vi.mocked(instancesApi.list).mockImplementation(() => 
