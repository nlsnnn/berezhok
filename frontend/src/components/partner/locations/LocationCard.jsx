import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { MapPin, Package, Tag, ChevronDown, ChevronUp, Check, Settings } from 'lucide-react'
import { useStores } from '@/context/StoresContext'
import LocationPinsSelector from './LocationPinsSelector'
import Spinner from '@/components/ui/feedback/Spinner'
import StatusBadge from '@/components/ui/feedback/StatusBadge'
import { LOCATION_STATUS } from '@/lib/constants'

function LocationCardBase({ location, boxCount = 0 }) {
  const { pinsStore } = useStores()
  const [pinsOpen, setPinsOpen] = useState(false)
  const [draft, setDraft] = useState(null) // null = not editing

  const currentPins = pinsStore.locationPins[location.id]
  const isSaving = pinsStore.savingFor === location.id

  useEffect(() => {
    if (pinsOpen) {
      if (pinsStore.availablePins.length === 0) {
        pinsStore.loadAvailable()
      }
      if (currentPins === undefined) {
        pinsStore.loadForLocation(location.id)
      }
    }
  }, [pinsOpen, currentPins, location.id, pinsStore])

  const handleEdit = () => {
    setDraft(currentPins?.map((p) => p.code) ?? [])
  }

  const handleSave = async () => {
    await pinsStore.save(location.id, draft)
    setDraft(null)
  }

  const handleCancel = () => {
    setDraft(null)
  }

  return (
    <article className="bg-white rounded-2xl border border-cream-200 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
      <div className="p-5">
        <div className="flex items-start gap-3">
          <div className="w-11 h-11 rounded-xl bg-brand-100 flex items-center justify-center shrink-0">
            <MapPin size={20} className="text-brand-600" />
          </div>

          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <h3 className="font-semibold text-brand-900 truncate">{location.name}</h3>
              <StatusBadge
                status={location.status}
                customLabel={LOCATION_STATUS[location.status]?.label}
                customColor={LOCATION_STATUS[location.status]?.color}
              />
            </div>
            <p className="text-sm text-brand-600 mt-1 line-clamp-2">{location.address}</p>

            <div className="mt-4 flex items-center justify-between">
              <div className="text-sm text-brand-500 flex items-center gap-1.5">
                <Package size={15} />
                <span>{boxCount} боксов</span>
              </div>
              <div className="flex items-center gap-3">
                <Link
                  to={`/partner/locations/${location.id}/edit`}
                  className="flex items-center gap-1 text-sm font-medium text-brand-600 hover:text-brand-800"
                >
                  <Settings size={14} />
                  Настроить
                </Link>
                <Link to="/partner/boxes" className="text-sm font-medium text-brand-600 hover:text-brand-800">
                  Смотреть боксы
                </Link>
              </div>
            </div>
          </div>
        </div>

        {/* Pins toggle */}
        <button
          type="button"
          onClick={() => setPinsOpen((v) => !v)}
          className="mt-4 w-full flex items-center justify-between text-sm text-brand-500 hover:text-brand-700 transition-colors"
        >
          <span className="flex items-center gap-1.5">
            <Tag size={14} />
            Пины заведения
            {currentPins && currentPins.length > 0 && (
              <span className="ml-1 text-xs bg-brand-100 text-brand-700 px-2 py-0.5 rounded-full">
                {currentPins.length}
              </span>
            )}
          </span>
          {pinsOpen ? <ChevronUp size={15} /> : <ChevronDown size={15} />}
        </button>
      </div>

      {pinsOpen && (
        <div className="border-t border-cream-100 px-5 pb-5 pt-4">
          {currentPins === undefined || pinsStore.loadingAvailable ? (
            <div className="flex justify-center py-4">
              <Spinner size={22} />
            </div>
          ) : draft !== null ? (
            /* Edit mode */
            <div className="space-y-3">
              <LocationPinsSelector
                availablePins={pinsStore.availablePins}
                selectedCodes={draft}
                onChange={setDraft}
              />
              <div className="flex gap-2 pt-1">
                <button
                  type="button"
                  onClick={handleSave}
                  disabled={isSaving}
                  className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-brand-600 text-white text-sm font-medium hover:bg-brand-700 disabled:opacity-50 transition-colors"
                >
                  {isSaving ? <Spinner size={14} /> : <Check size={14} />}
                  Сохранить
                </button>
                <button
                  type="button"
                  onClick={handleCancel}
                  className="px-3 py-1.5 rounded-lg border border-brand-200 text-brand-600 text-sm hover:border-brand-400 transition-colors"
                >
                  Отмена
                </button>
              </div>
            </div>
          ) : (
            /* View mode */
            <div>
              {currentPins.length === 0 ? (
                <p className="text-sm text-brand-400">Пины не выбраны</p>
              ) : (
                <div className="flex flex-wrap gap-1.5">
                  {currentPins.map((pin) => (
                    <span
                      key={pin.code}
                      className="px-2.5 py-1 rounded-full text-xs font-medium bg-brand-50 text-brand-700 border border-brand-100"
                    >
                      {pin.name_ru}
                    </span>
                  ))}
                </div>
              )}
              <button
                type="button"
                onClick={handleEdit}
                className="mt-3 text-sm text-brand-500 hover:text-brand-700 underline underline-offset-2 transition-colors"
              >
                Редактировать пины
              </button>
            </div>
          )}
        </div>
      )}
    </article>
  )
}

export default observer(LocationCardBase)
