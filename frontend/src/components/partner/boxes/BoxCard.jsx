import { Clock3, Pencil, Package, Trash2 } from 'lucide-react'
import { BOX_STATUS } from '@/lib/constants'
import { cn } from '@/lib/utils'

export default function BoxCard({ box, onEdit, onDelete }) {
  const status = BOX_STATUS[box.status] || { label: box.status, color: 'bg-gray-100 text-gray-800' }
  const hasDiscount = box.original_price && parseFloat(box.original_price) > parseFloat(box.discount_price)
  const discountPercent = hasDiscount
    ? Math.round((1 - parseFloat(box.discount_price) / parseFloat(box.original_price)) * 100)
    : 0

  return (
    <article className="bg-white rounded-2xl border border-cream-200 shadow-sm overflow-hidden hover:shadow-md transition-shadow">
      <div className="relative h-44 bg-cream-100">
        {box.image_url ? (
          <img src={box.image_url} alt={box.name} className="w-full h-full object-cover" />
        ) : (
          <div className="w-full h-full flex items-center justify-center">
            <Package size={42} className="text-cream-300" />
          </div>
        )}

        <div className="absolute top-3 left-3">
          <span className={cn('badge', status.color)}>{status.label}</span>
        </div>

        {hasDiscount && (
          <div className="absolute top-3 right-3">
            <span className="badge bg-red-500 text-white">-{discountPercent}%</span>
          </div>
        )}
      </div>

      <div className="p-4 space-y-3">
        <h3 className="font-semibold text-brand-900 line-clamp-1">{box.name}</h3>
        <p className="text-sm text-brand-600 line-clamp-2 min-h-[40px]">{box.description}</p>

        <div className="flex items-baseline gap-2">
          <span className="text-xl font-bold text-brand-700">{parseFloat(box.discount_price).toFixed(2)} ₽</span>
          {hasDiscount && <span className="text-sm text-cream-500 line-through">{parseFloat(box.original_price).toFixed(2)} ₽</span>}
        </div>

        <div className="flex items-center justify-between text-xs text-brand-500 pt-2 border-t border-cream-200">
          <div className="flex items-center gap-1">
            <Clock3 size={14} />
            <span>{box.pickup_time?.start} - {box.pickup_time?.end}</span>
          </div>
          <div className="flex items-center gap-1">
            <Package size={14} />
            <span>{box.quantity_available ?? box.quantity ?? 0} шт</span>
          </div>
        </div>

        <div className="flex gap-2 pt-1">
          <button className="btn-secondary w-full py-2 px-3 text-xs" onClick={() => onEdit(box)}>
            <Pencil size={14} />
            Редактировать
          </button>
          <button className="btn-danger w-full py-2 px-3 text-xs" onClick={() => onDelete(box)}>
            <Trash2 size={14} />
            Удалить
          </button>
        </div>
      </div>
    </article>
  )
}
