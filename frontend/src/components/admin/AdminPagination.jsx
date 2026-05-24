import Button from '@/components/ui/actions/Button'

export default function AdminPagination({ pagination, loading, onPrev, onNext }) {
  const total = pagination?.total ?? 0
  const limit = pagination?.limit ?? 20
  const offset = pagination?.offset ?? 0
  const shownTo = Math.min(total, offset + limit)

  return (
    <div className="flex flex-col gap-3 border-t border-cream-200 pt-4 sm:flex-row sm:items-center sm:justify-between">
      <p className="text-sm text-brand-600">
        Показано {total === 0 ? 0 : offset + 1}-{shownTo} из {total}
      </p>
      <div className="flex gap-2">
        <Button variant="secondary" onClick={onPrev} disabled={loading || offset <= 0}>
          Назад
        </Button>
        <Button variant="secondary" onClick={onNext} disabled={loading || !pagination?.has_more}>
          Вперёд
        </Button>
      </div>
    </div>
  )
}
