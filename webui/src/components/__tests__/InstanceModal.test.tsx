import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import InstanceDialog from '@/components/InstanceDialog'
import { BackendType, type Instance } from '@/types/instance'

// Mock the API module
vi.mock('@/lib/api', () => ({
  nodesApi: {
    list: vi.fn(() => Promise.resolve({})),
  },
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

describe('InstanceModal - Form Logic and Validation', () => {
  const mockOnSave = vi.fn()
  const mockOnOpenChange = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    window.sessionStorage.setItem('llamadeck_management_key', 'test-api-key-123')
    global.fetch = vi.fn(() => Promise.resolve(new Response(null, { status: 200 })))
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('Create Mode', () => {
    it('validates instance name is required', async () => {
      const user = userEvent.setup()

      render(
        <InstanceDialog
          open={true}
          onOpenChange={mockOnOpenChange}
          onSave={mockOnSave}
        />
      )

      // Navigate to the last tab where the save button is
      const advancedTab = screen.getByRole('tab', { name: /Advanced/i })
      await user.click(advancedTab)

      // Try to submit without name - save button should be disabled
      const saveButton = screen.getByTestId('dialog-save-button')
      expect(saveButton).toBeDisabled()

      // Go back to general tab to add name
      const generalTab = screen.getByRole('tab', { name: /General/i })
      await user.click(generalTab)

      // Add name
      const nameInput = screen.getByLabelText(/Instance Name/)
      await user.type(nameInput, 'test-instance')

      // Navigate back to advanced tab
      await user.click(advancedTab)

      await waitFor(() => {
        expect(screen.getByTestId('dialog-save-button')).not.toBeDisabled()
      })
    })

    it('validates instance name format', async () => {
      const user = userEvent.setup()

      render(
        <InstanceDialog
          open={true}
          onOpenChange={mockOnOpenChange}
          onSave={mockOnSave}
        />
      )

      const nameInput = screen.getByLabelText(/Instance Name/)

      // Test invalid characters
      await user.type(nameInput, 'test instance!')

      expect(screen.getByText(/can only contain letters, numbers, periods, hyphens, and underscores/)).toBeInTheDocument()

      // Navigate to advanced tab to check save button
      const advancedTab = screen.getByRole('tab', { name: /Advanced/i })
      await user.click(advancedTab)
      expect(screen.getByTestId('dialog-save-button')).toBeDisabled()

      // Go back to general tab to fix the name
      const generalTab = screen.getByRole('tab', { name: /General/i })
      await user.click(generalTab)

      // Select all and delete, then type valid name
      const newNameInput = screen.getByLabelText(/Instance Name/)
      await user.tripleClick(newNameInput)
      await user.keyboard('{Backspace}')
