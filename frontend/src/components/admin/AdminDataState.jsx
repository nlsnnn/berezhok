import Button from '@/components/ui/actions/Button'
import Spinner from '@/components/ui/feedback/Spinner'
import { getErrorMessage } from '@/lib/utils'

export default function AdminDataState({ loading, error, empty, emptyText, onRetry, children }) {
  if (loading) {
    return (
      <div className="card flex justify-center py-20">
        <Spinner size={34} />
      </div>
    )
  }

  if (error) {
    return (
      <div className="card flex flex-col items-center gap-4 py-14 text-center">
        <div>
          <h3 className="text-lg font-semibold text-brand-900">Не удалось загрузить данные</h3>
          <p className="mt-2 max-w-md text-sm text-brand-600">{getErrorMessage(error)}</p>
        </div>
        {onRetry && (
          <Button variant="secondary" onClick={onRetry}>
            Повторить
          </Button>
        )}
      </div>
    )
  }

  if (empty) {
    return (
      <div className="card py-14 text-center text-sm text-brand-600">
        {emptyText || 'Нет данных по текущим фильтрам'}
      </div>
    )
  }

  return children
}
