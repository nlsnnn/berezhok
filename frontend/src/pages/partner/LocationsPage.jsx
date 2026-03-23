import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { getPartnerProfile, listBoxes } from '@/api/partner'
import { Plus, MapPin } from 'lucide-react'
import PartnerNav from '@/components/PartnerNav'
import Spinner from '@/components/ui/Spinner'
import Button from '@/components/ui/Button'
import LocationCard from '@/components/ui/LocationCard'
import { useMemo } from 'react'

export default function LocationsPage() {
  const { data: profile, isLoading: isLoadingProfile, isError: isProfileError } = useQuery({
    queryKey: ['partner', 'profile'],
    queryFn: getPartnerProfile,
  })

  const { data: boxes, isLoading: isLoadingBoxes } = useQuery({
    queryKey: ['partner', 'boxes'],
    queryFn: listBoxes,
  })

  const locationBoxCounts = useMemo(() => {
    if (!boxes) return {}
    return boxes.reduce((acc, box) => {
      acc[box.location_id] = (acc[box.location_id] || 0) + 1
      return acc
    }, {})
  }, [boxes])

  const isLoading = isLoadingProfile || isLoadingBoxes

  return (
    <div className="min-h-screen flex flex-col bg-cream-50">
      <PartnerNav />

      <main className="flex-1 max-w-5xl mx-auto w-full px-4 sm:px-6 py-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-3xl font-bold text-brand-900">Мои локации</h1>
            <p className="text-brand-600 mt-1">
              Управляйте точками продаж
            </p>
          </div>
          <Link to="/partner/locations/new">
            <Button className="gap-2">
              <Plus size={18} />
              Добавить локацию
            </Button>
          </Link>
        </div>

        {/* Loading */}
        {isLoading && (
          <div className="flex justify-center py-20">
            <Spinner size={32} />
          </div>
        )}

        {/* Error */}
        {isProfileError && (
          <div className="card text-center py-12 text-red-500">
            Не удалось загрузить локации
          </div>
        )}

        {/* Empty state */}
        {profile && (!profile.locations || profile.locations.length === 0) && (
          <div className="card text-center py-16">
            <div className="w-20 h-20 rounded-full bg-cream-200 flex items-center justify-center mx-auto mb-4">
              <MapPin size={40} className="text-cream-400" />
            </div>
            <h3 className="font-semibold text-brand-700 text-lg mb-2">
              Нет локаций
            </h3>
            <p className="text-brand-500 mb-6">
              Добавьте первую точку продаж
            </p>
            <Link to="/partner/locations/new">
              <Button className="gap-2">
                <Plus size={18} />
                Добавить локацию
              </Button>
            </Link>
          </div>
        )}

        {/* Locations grid */}
        {profile && profile.locations && profile.locations.length > 0 && (
          <div className="grid md:grid-cols-2 gap-6">
            {profile.locations.map((location) => (
              <LocationCard
                key={location.id}
                location={location}
                boxCount={locationBoxCounts[location.id] || 0}
              />
            ))}
          </div>
        )}
      </main>
    </div>
  )
}
