import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { AuthProvider } from '@/contexts/AuthContext'
import LoginDialog from '@/components/LoginDialog'
import { authApi } from '@/lib/api'

vi.mock('@/lib/api', () => ({
  authApi: {
    whoami: vi.fn(),
  },
}))

const whoamiMock = vi.mocked(authApi.whoami)

function renderLogin() {
  return render(
    <AuthProvider>
      <LoginDialog open />
    </AuthProvider>
