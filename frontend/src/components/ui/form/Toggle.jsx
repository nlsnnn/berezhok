import { cn } from '@/lib/utils'

export default function Toggle({
  id,
  checked,
  onChange,
  disabled = false,
  label,
  description,
  className,
}) {
  const handleClick = () => {
    if (disabled) return
    onChange?.(!checked)
  }

  return (
    <div className={cn('flex items-start gap-3', className)}>
      <button
        type="button"
        role="switch"
        id={id}
        aria-checked={checked}
        disabled={disabled}
        onClick={handleClick}
        className={cn(
          'relative inline-flex h-6 w-11 shrink-0 items-center rounded-full transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-brand-400 focus:ring-offset-2 cursor-pointer',
          checked ? 'bg-brand-500' : 'bg-cream-300',
          disabled && 'opacity-50 cursor-not-allowed',
        )}
      >
        <span
          className={cn(
            'inline-block h-5 w-5 transform rounded-full bg-white shadow transition-transform duration-200',
            checked ? 'translate-x-5' : 'translate-x-0.5',
          )}
        />
      </button>
      {(label || description) && (
        <div className="flex-1 min-w-0">
          {label && (
            <label
              htmlFor={id}
              onClick={handleClick}
              className={cn(
                'block text-sm font-medium text-brand-700',
                !disabled && 'cursor-pointer',
              )}
            >
              {label}
            </label>
          )}
          {description && <p className="text-xs text-brand-500 mt-0.5">{description}</p>}
        </div>
      )}
    </div>
  )
}
