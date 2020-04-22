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
