/**
 * Parses error response from API calls and returns a formatted error message
 */
export async function parseErrorResponse(response: Response): Promise<string> {
  let errorMessage = `HTTP ${response.status}`

  try {
    const errorText = await response.text()
