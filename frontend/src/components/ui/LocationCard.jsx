import { MapPin, Package } from 'lucide-react'
import { Link } from 'react-router-dom'

export default function LocationCard({ location, boxCount = 0 }) {
  return (
    <div className="card hover:shadow-md transition-shadow">
      <div className="flex items-start gap-4">
        {/* Icon */}
        <div className="w-12 h-12 rounded-xl bg-brand-100 flex items-center justify-center shrink-0">
          <MapPin size={24} className="text-brand-500" />
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0">
          <h3 className="font-semibold text-brand-900 mb-1 truncate">
            {location.name}
          </h3>
          <p className="text-sm text-brand-600 mb-3 line-clamp-2">
            {location.address}
          </p>

          <div className="flex items-center justify-between gap-4">
            {/* Box count */}
            <div className="flex items-center gap-1.5 text-sm text-brand-500">
              <Package size={16} />
              <span>{boxCount} {boxCount === 1 ? 'бокс' : boxCount < 5 ? 'бокса' : 'боксов'}</span>
            </div>

            {/* Link */}
            <Link
              to="/partner/boxes"
              className="text-sm text-brand-500 hover:text-brand-700 font-medium transition-colors"
            >
              Смотреть боксы →
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}
