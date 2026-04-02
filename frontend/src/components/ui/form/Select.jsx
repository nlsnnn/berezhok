import { cn } from '@/lib/utils'

export default function Select({ className, error, children, ...props }) {
  return (
    <div className="w-full">
      <select
        className={cn('input-base appearance-none bg-white cursor-pointer', error && 'border-red-400 focus:ring-red-400', className)}
        {...props}
      >
        {children}
      </select>
      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}
    </div>
  )
}
