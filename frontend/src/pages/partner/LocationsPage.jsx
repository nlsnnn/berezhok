import { useEffect, useMemo } from 'react'
import { observer } from 'mobx-react-lite'
import { Link } from 'react-router-dom'
import { Plus } from 'lucide-react'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import LocationCard from '@/components/partner/locations/LocationCard'
import Spinner from '@/components/ui/feedback/Spinner'
import Button from '@/components/ui/actions/Button'
import { useStores } from '@/context/StoresContext'

function LocationsPageBase() {
  const { locationsStore, boxesStore } = useStores()

  useEffect(() => {
    locationsStore.loadProfile()
    boxesStore.load()
  }, [boxesStore, locationsStore])

  const counts = useMemo(() => {
    return boxesStore.items.reduce((acc, box) => {
      const locationId = box.location_id
      if (!locationId) return acc
      acc[locationId] = (acc[locationId] || 0) + 1
      return acc
    }, {})
  }, [boxesStore.items])

  return (
    <PartnerLayout
      title="Локации"
      subtitle="Точки продаж и их активность"
      actions={
        <Link to="/partner/locations/new">
          <Button className="gap-2">
            <Plus size={16} />
            Добавить локацию
          </Button>
        </Link>
      }
    >
      {(locationsStore.loading || boxesStore.loading) && (
        <div className="flex justify-center py-24">
          <Spinner size={34} />
        </div>
      )}

      {locationsStore.error && (
        <div className="card text-center py-12 text-red-600">Не удалось загрузить локации</div>
      )}

      {!locationsStore.loading && !locationsStore.error && locationsStore.locations.length === 0 && (
        <div className="card text-center py-14">
          <p className="text-brand-700 font-medium">Нет локаций</p>
          <p className="text-sm text-brand-500 mt-1">Добавьте первую точку продаж</p>
          <Link to="/partner/locations/new">
            <Button className="mt-5">Добавить</Button>
          </Link>
        </div>
      )}

      {!locationsStore.loading && locationsStore.locations.length > 0 && (
        <div className="grid lg:grid-cols-2 gap-5">
          {locationsStore.locations.map((location) => (
            <LocationCard key={location.id} location={location} boxCount={counts[location.id] || 0} />
          ))}
        </div>
      )}
    </PartnerLayout>
  )
}

export default observer(LocationsPageBase)
