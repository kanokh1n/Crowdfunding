import { useToastStore } from '@/stores/toast'
import { ApiError } from '@/api/client'

export function useErrorToast() {
  const toast = useToastStore()

  function showError(err: unknown) {
    if (err instanceof ApiError) {
      toast.add(err.userMessage, 'error')
    } else if (err instanceof Error) {
      toast.add(err.message || 'Произошла ошибка.', 'error')
    } else {
      toast.add('Произошла непредвиденная ошибка.', 'error')
    }
  }

  function showSuccess(message: string) {
    toast.add(message, 'success')
  }

  return { showError, showSuccess }
}
