import { forwardRef } from 'react'
import Input from './Input'

export function formatPhoneDisplay(raw = '') {
  const digits = String(raw).replace(/\D/g, '')
  let d = digits
  if (d.startsWith('8')) d = '7' + d.slice(1)
  if (d.length > 0 && !d.startsWith('7')) d = '7' + d
  d = d.slice(0, 11)
  if (!d) return ''
  const r = d.slice(1)
  let out = '+7'
  if (r.length > 0) out += ' (' + r.slice(0, Math.min(3, r.length))
  if (r.length >= 3) out += ')'
  if (r.length > 3) out += ' ' + r.slice(3, Math.min(6, r.length))
  if (r.length > 6) out += '-' + r.slice(6, Math.min(8, r.length))
  if (r.length > 8) out += '-' + r.slice(8, Math.min(10, r.length))
  return out
}

export function parsePhoneInput(input = '') {
  const digits = String(input).replace(/\D/g, '')
  let d = digits
  if (d.startsWith('8')) d = '7' + d.slice(1)
  if (d.length > 0 && !d.startsWith('7')) d = '7' + d
  d = d.slice(0, 11)
  return d.length > 0 ? '+' + d : ''
}

const PhoneInput = forwardRef(function PhoneInput({ value = '', onChange, ...props }, ref) {
  const handleChange = (e) => {
    const raw = parsePhoneInput(e.target.value)
    onChange?.({ target: { value: raw } })
  }

  return (
    <Input
      ref={ref}
      type="tel"
      value={formatPhoneDisplay(value)}
      onChange={handleChange}
      placeholder="+7 (XXX) XXX-XX-XX"
      {...props}
    />
  )
})

export default PhoneInput
