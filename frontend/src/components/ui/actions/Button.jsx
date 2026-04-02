import { cn } from '@/lib/utils'

export default function Button({ className, variant = 'primary', size = 'md', children, ...props }) {
  const base = {
    primary: 'btn-primary',
    secondary: 'btn-secondary',
    danger: 'btn-danger',
    ghost: 'btn-ghost',
  }[variant]

  const sizes = {
    sm: 'px-4 py-1.5 text-xs',
    md: '',
    lg: 'px-8 py-3 text-base',
  }[size]

  return (
    <button className={cn(base, sizes, className)} {...props}>
      {children}
    </button>
  )
}
