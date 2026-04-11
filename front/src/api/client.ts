const API_BASE_URL = '/api'

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
}

interface FetchOptions extends RequestInit {
  requireAuth?: boolean
}

async function apiFetch<T>(endpoint: string, options: FetchOptions = {}): Promise<T> {
  const { requireAuth = false, headers: customHeaders, ...rest } = options

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(customHeaders as Record<string, string>),
  }

  if (requireAuth) {
    const token = localStorage.getItem('access_token')
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    headers,
    ...rest,
  })

  if (!response.ok) {
    const errorBody = await response.json().catch(() => null)

    if (response.status === 401 && requireAuth) {
      try {
        await refreshAccessToken()
        const newToken = localStorage.getItem('access_token')
        if (newToken) {
          headers['Authorization'] = `Bearer ${newToken}`
        }
        const retryResponse = await fetch(`${API_BASE_URL}${endpoint}`, {
          headers,
          ...rest,
        })
        if (!retryResponse.ok) {
          throw new ApiError(retryResponse.status, await retryResponse.json().catch(() => null))
        }
        return retryResponse.json() as Promise<T>
      } catch {
        logoutAndRedirect()
        throw new ApiError(401, errorBody)
      }
    }

    throw new ApiError(response.status, errorBody)
  }

  if (response.status === 204) {
    return undefined as unknown as T
  }

  return response.json() as Promise<T>
}

/** Fetch paginated list and return just the data array */
async function fetchList<T>(endpoint: string, options?: FetchOptions): Promise<T[]> {
  const resp = await apiFetch<PaginatedResponse<T>>(endpoint, options)
  return resp.data
}

async function refreshAccessToken(): Promise<void> {
  const refreshToken = localStorage.getItem('refresh_token')
  if (!refreshToken) throw new Error('No refresh token')

  const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken }),
  })

  if (!response.ok) throw new Error('Refresh failed')

  const data = await response.json()
  localStorage.setItem('access_token', data.access_token)
  // Note: backend does NOT return a new refresh_token, keep the existing one
}

function logoutAndRedirect() {
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
  window.location.href = '/auth'
}

class ApiError extends Error {
  constructor(
    public status: number,
    public data?: unknown
  ) {
    super(`API Error: ${status}`)
    this.name = 'ApiError'
  }
}

export { apiFetch, fetchList, ApiError }
