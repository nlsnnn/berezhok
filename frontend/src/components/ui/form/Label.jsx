import { cn } from '@/lib/utils'

export default function Label({ children, required, className, ...props }) {
  return (
    <label className={cn('block text-sm font-medium text-brand-700 mb-1', className)} {...props}>
      {children}
      {required && <span className="text-red-400 ml-1">*</span>}
    </label>
  )
}
