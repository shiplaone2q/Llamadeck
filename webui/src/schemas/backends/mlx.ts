import { z } from 'zod'

// Define the MLX backend options schema
export const MlxBackendOptionsSchema = z.object({
  // Basic connection options
  model: z.string().optional(),
  host: z.string().optional(),
  port: z.number().optional(),

  // Model and adapter options
  adapter_path: z.string().optional(),
  draft_model: z.string().optional(),
  num_draft_tokens: z.number().optional(),
  trust_remote_code: z.boolean().optional(),

  // Logging and templates
  log_level: z.enum(['DEBUG', 'INFO', 'WARNING', 'ERROR', 'CRITICAL']).optional(),
