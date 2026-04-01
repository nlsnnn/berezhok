import { Link } from 'react-router-dom'
import { MapPin, Package } from 'lucide-react'

export default function LocationCard({ location, boxCount = 0 }) {
  return (
    <article className="bg-white rounded-2xl border border-cream-200 p-5 shadow-sm hover:shadow-md transition-shadow">
      <div className="flex items-start gap-3">
        <div className="w-11 h-11 rounded-xl bg-brand-100 flex items-center justify-center shrink-0">
          <MapPin size={20} className="text-brand-600" />
        </div>

        <div className="flex-1 min-w-0">
          <h3 className="font-semibold text-brand-900 truncate">{location.name}</h3>
          <p className="text-sm text-brand-600 mt-1 line-clamp-2">{location.address}</p>

          <div className="mt-4 flex items-center justify-between">
            <div className="text-sm text-brand-500 flex items-center gap-1.5">
              <Package size={15} />
              <span>{boxCount} боксов</span>
            </div>
            <Link to="/partner/boxes" className="text-sm font-medium text-brand-600 hover:text-brand-800">
              Смотреть боксы
            </Link>
          </div>
        </div>
      </div>
    </article>
  )
}
