import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { createBox, getPartnerProfile } from '@/api/partner'
import { toast } from 'sonner'
import { ArrowLeft } from 'lucide-react'
import PartnerNav from '@/components/PartnerNav'
import BoxForm from '@/components/BoxForm'
import Spinner from '@/components/ui/Spinner'

export default function CreateBoxPage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data: profile, isLoading: isLoadingProfile } = useQuery({
    queryKey: ['partner', 'profile'],
    queryFn: getPartnerProfile,
  })

  const createMutation = useMutation({
    mutationFn: createBox,
    onSuccess: () => {
      queryClient.invalidateQueries(['partner', 'boxes'])
      toast.success('Бокс создан')
      navigate('/partner/boxes')
    },
    onError: (error) => {
      console.error('Create error:', error)
      const message = error.response?.data?.message || 'Не удалось создать бокс'
      toast.error(message)
    },
  })

  const handleSubmit = (formData) => {
    // Transform form data to match backend expectations
    const payload = {
      location_id: formData.location_id,
      name: formData.name,
      description: formData.description,
      discount_price: formData.discount_price,
      original_price: formData.original_price || null,
      pickup_time_start: formData.pickup_time_start,
      pickup_time_end: formData.pickup_time_end,
      quantity: parseInt(formData.quantity),
      image_url: formData.image_url || '',
      status: formData.status,
    }
    createMutation.mutate(payload)
  }

  return (
    <div className="min-h-screen flex flex-col bg-cream-50">
      <PartnerNav />

      <main className="flex-1 max-w-3xl mx-auto w-full px-4 sm:px-6 py-8">
        {/* Header */}
        <div className="mb-8">
          <button
            onClick={() => navigate('/partner/boxes')}
            className="flex items-center gap-2 text-brand-600 hover:text-brand-800 mb-4 transition-colors"
          >
            <ArrowLeft size={18} />
            Назад к боксам
          </button>
          <h1 className="text-3xl font-bold text-brand-900">Создать бокс</h1>
          <p className="text-brand-600 mt-1">
            Заполните информацию о новом предложении
          </p>
        </div>

        {/* Form */}
        <div className="card">
          {isLoadingProfile ? (
            <div className="flex justify-center py-12">
              <Spinner size={32} />
            </div>
          ) : (
            <BoxForm
              locations={profile?.locations || []}
              onSubmit={handleSubmit}
              isLoading={createMutation.isPending}
            />
          )}
        </div>
      </main>
    </div>
  )
}
