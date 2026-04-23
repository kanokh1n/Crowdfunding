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

const STATUS_MESSAGES: Record<number, string> = {
  400: 'Некорректный запрос. Проверьте введённые данные.',
  401: 'Сессия истекла. Пожалуйста, войдите снова.',
  403: 'У вас нет прав для этого действия.',
  404: 'Запрашиваемый ресурс не найден.',
  409: 'Конфликт данных. Возможно, такая запись уже существует.',
  413: 'Файл слишком большой.',
  422: 'Ошибка валидации. Проверьте правильность данных.',
  429: 'Слишком много запросов. Подождите немного.',
  500: 'Внутренняя ошибка сервера. Попробуйте позже.',
  502: 'Сервер временно недоступен.',
  503: 'Сервис недоступен. Попробуйте позже.',
}

const BACKEND_MESSAGES: Record<string, string> = {
  'email already exists':                    'Пользователь с таким email уже зарегистрирован.',
  'invalid credentials':                     'Неверный email или пароль.',
  'invalid token':                           'Недействительный токен. Войдите снова.',
  'token expired':                           'Токен истёк. Войдите снова.',
  'project not found':                       'Проект не найден.',
  'user not found':                          'Пользователь не найден.',
  'comment not found':                       'Комментарий не найден.',
  'forbidden':                               'У вас нет прав для этого действия.',
  'unauthorized':                            'Необходима авторизация.',
  'file required':                           'Файл не выбран.',
  'file too large (max 5mb)':               'Файл слишком большой. Максимум 5 МБ.',
  'file too large (max 20mb)':              'Файл слишком большой. Максимум 20 МБ.',
  'unsupported format, use jpg/png/webp/gif': 'Неподдерживаемый формат. Используйте JPG, PNG, WebP или GIF.',
  'unsupported format, use pdf or pptx':    'Неподдерживаемый формат. Используйте PDF или PPTX.',
  'failed to create project':               'Не удалось создать проект. Попробуйте позже.',
  'failed to fetch projects':               'Не удалось загрузить проекты.',
  'failed to save file':                    'Не удалось сохранить файл.',
  'storage error':                          'Ошибка хранилища файлов.',
  'already liked':                          'Вы уже поддержали этот проект.',
  'not liked':                              'Вы ещё не поддерживали этот проект.',
  'insufficient funds':                     'Недостаточно средств.',
  'invalid amount':                         'Некорректная сумма.',
  'project is not active':                  'Проект не активен.',
  'email not verified':                     'Email не подтверждён. Проверьте почту.',
}

class ApiError extends Error {
  constructor(
    public status: number,
    public data?: unknown
  ) {
    super(`API Error: ${status}`)
    this.name = 'ApiError'
  }

  get userMessage(): string {
    const raw = (this.data as Record<string, string> | null)?.error
    if (raw) {
      const key = raw.toLowerCase().trim()
      if (BACKEND_MESSAGES[key]) return BACKEND_MESSAGES[key]
      // Частичное совпадение для длинных сообщений с деталями
      for (const [pattern, msg] of Object.entries(BACKEND_MESSAGES)) {
        if (key.includes(pattern)) return msg
      }
      // Если бэкенд прислал читаемое сообщение — показываем его
      if (raw.length < 120) return raw
    }
    return STATUS_MESSAGES[this.status] ?? 'Произошла непредвиденная ошибка. Попробуйте позже.'
  }
}

export { apiFetch, fetchList, ApiError }
