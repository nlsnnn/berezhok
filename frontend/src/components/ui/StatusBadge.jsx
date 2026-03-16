import { APPLICATION_STATUS } from '@/lib/constants'
import { cn } from '@/lib/utils'

export default function StatusBadge({ status, customLabel, customColor }) {
  const cfg = APPLICATION_STATUS[status]
  const label = customLabel ?? cfg?.label ?? status
  const color = customColor ?? cfg?.color ?? 'bg-gray-100 text-gray-700'

  return (
    <span className={cn('badge', color)}>
      {label}
    </span>
  )
}
