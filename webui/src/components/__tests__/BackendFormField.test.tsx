import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import BackendFormField from '@/components/BackendFormField'
import type { CreateInstanceOptions } from '@/schemas/instanceOptions'

describe('BackendFormField - models_preset Field States', () => {
  const mockOnChange = vi.fn()

  describe('models_preset field with no preset_ini content', () => {
    it('shows default state when models_preset is empty and no preset_ini', () => {
      render(
        <BackendFormField
          fieldKey="models_preset"
          value=""
          onChange={mockOnChange}
          formData={{}}
        />
      )

      expect(screen.getByLabelText('Models Preset Path')).toBeInTheDocument()
      expect(screen.getByPlaceholderText('/path/to/preset.ini')).toBeInTheDocument()
      expect(screen.getByText(/Optional: Path to preset.ini for router mode/)).toBeInTheDocument()
      expect(screen.queryByText('Auto')).not.toBeInTheDocument()
      expect(screen.queryByText('Custom')).not.toBeInTheDocument()
    })

    it('shows Auto badge when models_preset is empty but preset_ini has content', () => {
      const formData: CreateInstanceOptions = {
        preset_ini: '[model1]\nmodel = /path/to/model.gguf\n'
      }

      render(
        <BackendFormField
          fieldKey="models_preset"
          value=""
          onChange={mockOnChange}
          formData={formData}
        />
      )

      expect(screen.getByLabelText('Models Preset Path')).toBeInTheDocument()
      expect(screen.getByText('Auto')).toBeInTheDocument()
      expect(screen.getByRole('alert')).toBeInTheDocument()
      expect(screen.getByText('Preset defined')).toBeInTheDocument()
      expect(screen.getByText('Will be auto-set to the preset.ini created in Preset Editor')).toBeInTheDocument()
    })

    it('shows plain help text instead of callout when custom path is set', () => {
      const formData: CreateInstanceOptions = {
        preset_ini: '[model1]\nmodel = /path/to/model.gguf\n'
      }

      render(
        <BackendFormField
          fieldKey="models_preset"
          value="/custom/path/to/preset.ini"
          onChange={mockOnChange}
          formData={formData}
        />
      )

      expect(screen.queryByRole('alert')).not.toBeInTheDocument()
      expect(screen.getByText(/Optional: Path to preset.ini for router mode/)).toBeInTheDocument()
    })
  })

  describe('models_preset field with custom path', () => {
    it('shows Custom badge when models_preset has user-provided value', () => {
      const formData: CreateInstanceOptions = {
        preset_ini: '[model1]\nmodel = /path/to/model.gguf\n'
      }

      render(
        <BackendFormField
          fieldKey="models_preset"
          value="/custom/path/to/preset.ini"
          onChange={mockOnChange}
          formData={formData}
        />
