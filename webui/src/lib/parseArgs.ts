// Splits a string into arguments the way a shell would: whitespace
// separates arguments, single and double quotes group values containing
// spaces or commas (e.g. --tensor-split "3,1").
export function parseArgs(input: string): string[] {
  const args: string[] = []
  let current = ''
  let quote: string | null = null
  let hasToken = false

  for (const ch of input) {
    if (quote) {
      if (ch === quote) {
        quote = null
      } else {
        current += ch
      }
