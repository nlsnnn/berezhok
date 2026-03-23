import { Pencil, Trash2, Clock, Package } from 'lucide-react'
import { BOX_STATUS } from '@/lib/constants'
import { cn } from '@/lib/utils'

export default function BoxCard({ box, onEdit, onDelete }) {
  const status = BOX_STATUS[box.status] || { label: box.status, color: 'bg-gray-100 text-gray-800' }
  
  const hasDiscount = box.original_price && parseFloat(box.original_price) > parseFloat(box.discount_price)
  const discountPercent = hasDiscount 
    ? Math.round((1 - parseFloat(box.discount_price) / parseFloat(box.original_price)) * 100)
    : 0

  return (
    <div className="card group hover:shadow-lg transition-shadow">
      {/* Image */}
      <div className="relative w-full h-48 bg-cream-100 rounded-t-xl overflow-hidden mb-4">
        {box.image_url ? (
          <img
            src={box.image_url}
            alt={box.name}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center">
            <Package size={48} className="text-cream-300" />
          </div>
        )}
        
        {/* Status badge */}
        <div className="absolute top-2 left-2">
          <span className={cn('badge', status.color)}>{status.label}</span>
        </div>

        {/* Discount badge */}
        {hasDiscount && (
          <div className="absolute top-2 right-2">
            <span className="badge bg-red-500 text-white">-{discountPercent}%</span>
          </div>
        )}

        {/* Actions - show on hover */}
        <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-40 transition-all flex items-center justify-center gap-2 opacity-0 group-hover:opacity-100">
          <button
            onClick={() => onEdit(box)}
            className="p-2 bg-white rounded-lg hover:bg-brand-500 hover:text-white transition-colors"
            title="Редактировать"
          >
            <Pencil size={18} />
          </button>
          <button
            onClick={() => onDelete(box)}
            className="p-2 bg-white rounded-lg hover:bg-red-500 hover:text-white transition-colors"
            title="Удалить"
          >
            <Trash2 size={18} />
          </button>
        </div>
      </div>

      {/* Content */}
      <div className="space-y-3">
        {/* Title */}
        <h3 className="font-semibold text-brand-900 text-lg line-clamp-1">
          {box.name}
        </h3>

        {/* Description */}
        <p className="text-sm text-brand-600 line-clamp-2">
          {box.description}
        </p>

        {/* Price */}
        <div className="flex items-baseline gap-2">
          <span className="text-2xl font-bold text-green-600">
            {parseFloat(box.discount_price).toFixed(2)} ₽
          </span>
          {hasDiscount && (
            <span className="text-sm text-gray-500 line-through">
              {parseFloat(box.original_price).toFixed(2)} ₽
            </span>
          )}
        </div>

        {/* Meta info */}
        <div className="flex items-center justify-between text-xs text-brand-500 pt-2 border-t border-cream-200">
          <div className="flex items-center gap-1">
            <Clock size={14} />
            <span>{box.pickup_time?.start} - {box.pickup_time?.end}</span>
          </div>
          <div className="flex items-center gap-1">
            <Package size={14} />
            <span>{box.quantity} шт</span>
          </div>
        </div>
      </div>
    </div>
  )
}
