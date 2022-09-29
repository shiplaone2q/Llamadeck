import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import InstanceCard from '@/components/InstanceCard'
import { type Instance, BackendType } from '@/types/instance'

// Mock the health hook since we're not testing health logic here
vi.mock('@/hooks/useInstanceHealth', () => ({
  useInstanceHealth: vi.fn(() => ({
    state: 'ready',
    instanceStatus: 'running',
    lastChecked: new Date(),
    source: 'http'
  }))
}))

describe('InstanceCard - Instance Actions and State', () => {
  const mockStartInstance = vi.fn()
  const mockStopInstance = vi.fn()
  const mockDeleteInstance = vi.fn()
  const mockEditInstance = vi.fn()

  const stoppedInstance: Instance = {
    id: 1,
    name: 'test-instance',
    status: 'stopped',
    options: { backend_type: BackendType.LLAMA_CPP, backend_options: { model: 'test-model.gguf' } }
  }

  const runningInstance: Instance = {
    id: 2,
    name: 'running-instance',
    status: 'running',
    options: { backend_type: BackendType.LLAMA_CPP, backend_options: { model: 'running-model.gguf' } }
  }

beforeEach(() => {
  vi.clearAllMocks()
  window.sessionStorage.setItem('llamadeck_management_key', 'test-api-key-123')
  global.fetch = vi.fn(() => Promise.resolve(new Response(null, { status: 200 })))
})

afterEach(() => {
  vi.restoreAllMocks()
})

  describe('Instance Action Buttons', () => {
    it('calls startInstance when start button clicked on stopped instance', async () => {
      const user = userEvent.setup()
      
      render(
        <InstanceCard
          instance={stoppedInstance}
          startInstance={mockStartInstance}
          stopInstance={mockStopInstance}
          deleteInstance={mockDeleteInstance}
          editInstance={mockEditInstance}
        />
      )

      const startButton = screen.getByTitle('Start instance')
      expect(startButton).not.toBeDisabled()
      
      await user.click(startButton)
      
      expect(mockStartInstance).toHaveBeenCalledWith('test-instance')
    })

    it('calls stopInstance when stop button clicked on running instance', async () => {
      const user = userEvent.setup()
      
      render(
        <InstanceCard
          instance={runningInstance}
          startInstance={mockStartInstance}
          stopInstance={mockStopInstance}
          deleteInstance={mockDeleteInstance}
          editInstance={mockEditInstance}
        />
      )

      const stopButton = screen.getByTitle('Stop instance')
      expect(stopButton).not.toBeDisabled()
      
      await user.click(stopButton)
      
      expect(mockStopInstance).toHaveBeenCalledWith('running-instance')
    })

    it('calls editInstance when edit button clicked', async () => {
      const user = userEvent.setup()
      
      render(
        <InstanceCard
          instance={stoppedInstance}
          startInstance={mockStartInstance}
          stopInstance={mockStopInstance}
          deleteInstance={mockDeleteInstance}
          editInstance={mockEditInstance}
        />
      )

      const editButton = screen.getByTitle('Edit instance')
      await user.click(editButton)
      
      expect(mockEditInstance).toHaveBeenCalledWith(stoppedInstance)
    })

    it('opens logs dialog when logs button clicked', async () => {
      const user = userEvent.setup()

      render(
        <InstanceCard
          instance={stoppedInstance}
          startInstance={mockStartInstance}
          stopInstance={mockStopInstance}
          deleteInstance={mockDeleteInstance}
          editInstance={mockEditInstance}
        />
      )

      // First click "More actions" to reveal the logs button
      const moreActionsButton = screen.getByTitle('More actions')
      await user.click(moreActionsButton)

      const logsButton = screen.getByTitle('View logs')
      await user.click(logsButton)

      // Should open logs dialog (we can verify this by checking if dialog title appears)
      expect(screen.getByText(`Logs: ${stoppedInstance.name}`)).toBeInTheDocument()
    })
  })

