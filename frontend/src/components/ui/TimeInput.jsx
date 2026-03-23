import { cn } from '@/lib/utils'
import { forwardRef } from 'react'

const TimeInput = forwardRef(function TimeInput({ className, error, label, ...props }, ref) {
  return (
    <div className="w-full">
      {label && <label className="block text-sm font-medium text-brand-700 mb-1">{label}</label>}
      <input
        ref={ref}
        type="time"
        className={cn(
          'input-base',
          error && 'border-red-400 focus:ring-red-400',
          className
        )}
        {...props}
      />
      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}
    </div>
  )
})

export default TimeInput
