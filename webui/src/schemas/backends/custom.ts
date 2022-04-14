import { z } from 'zod'

// Custom backend options reference a named backends.custom.<name> config
// entry. Args are literal strings; {port} and {model} are substituted by
// the server before launch.
export const CustomBackendOptionsSchema = z.object({
  name: z.string(),
  host: z.string().optional(),
  port: z.number().optional(),
  model: z.string().optional(),
  args: z.array(z.string()).optional(),
})

// Infer the TypeScript type from the schema
export type CustomBackendOptions = z.infer<typeof CustomBackendOptionsSchema>
