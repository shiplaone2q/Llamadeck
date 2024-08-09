import { BackendType } from '@/types/instance'
import { z } from 'zod'

// Import backend schemas from separate files
import {
  LlamaCppBackendOptionsSchema,
  type LlamaCppBackendOptions,
  getAllLlamaCppFieldKeys,
  getLlamaCppFieldType,
  getAllLlamaCppAltKeys,
  getLlamaCppAltKeyType,
  MlxBackendOptionsSchema,
  type MlxBackendOptions,
  getAllMlxFieldKeys,
  getMlxFieldType,
  VllmBackendOptionsSchema,
  type VllmBackendOptions,
  getAllVllmFieldKeys,
  getVllmFieldType,
  CustomBackendOptionsSchema,
  type CustomBackendOptions,
  getAllCustomFieldKeys,
  getCustomFieldType
} from './backends'

// Backend options union. Custom is first: its required `name` field
// discriminates it from the all-optional built-in option objects.
export const BackendOptionsSchema = z.union([
  CustomBackendOptionsSchema,
  LlamaCppBackendOptionsSchema,
  MlxBackendOptionsSchema,
  VllmBackendOptionsSchema,
])

// Define the main create instance options schema
export const CreateInstanceOptionsSchema = z.object({
  // Restart options
  auto_restart: z.boolean().optional(),
  max_restarts: z.number().optional(),
  restart_delay: z.number().optional(),
  idle_timeout: z.number().optional(),
  on_demand_start: z.boolean().optional(),

  // Environment variables
  environment: z.record(z.string(), z.string()).optional(),
