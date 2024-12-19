import { instancesApi } from '@/lib/api'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock fetch globally
const mockFetch = vi.fn()

describe('API Error Handling', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Set the mock fetch in beforeEach to ensure it overrides the setup.ts mock
    global.fetch = mockFetch
  })

  it('converts HTTP errors to meaningful messages', async () => {
    const mockResponse = new Response('Instance already exists', {
      status: 409,
      statusText: 'Conflict'
    })
    mockFetch.mockResolvedValue(mockResponse)

    await expect(instancesApi.create('existing', {}))
      .rejects
      .toThrow('HTTP 409: Instance already exists')
  })

  it('handles empty error responses gracefully', async () => {
    const mockResponse = new Response('', {
      status: 500,
      statusText: 'Internal Server Error'
    })
    mockFetch.mockResolvedValue(mockResponse)

    await expect(instancesApi.list())
      .rejects
