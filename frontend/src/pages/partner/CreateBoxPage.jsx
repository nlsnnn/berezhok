import { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { useNavigate } from 'react-router-dom'
import { ArrowLeft } from 'lucide-react'
import { toast } from 'sonner'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import BoxForm from '@/components/partner/boxes/BoxForm'
import Spinner from '@/components/ui/feedback/Spinner'
import Button from '@/components/ui/actions/Button'
import { useStores } from '@/context/StoresContext'

function CreateBoxPageBase() {
  const navigate = useNavigate()
  const { boxesStore, locationsStore } = useStores()

  useEffect(() => {
    locationsStore.loadProfile()
  }, [locationsStore])

  const handleSubmit = async (formData) => {
    const payload = {
      location_id: formData.location_id,
      name: formData.name,
      description: formData.description,
      original_price: formData.original_price || null,
      discount_price: Number(formData.discount_price),
      quantity: Number(formData.quantity_available),
      pickup_time_start: formData.pickup_time_start,
      pickup_time_end: formData.pickup_time_end,
      image_url: formData.image_url || '',
      status: formData.status,
    }

    try {
      await boxesStore.create(payload)
      toast.success('Бокс создан')
      navigate('/partner/boxes')
    } catch {
      toast.error('Не удалось создать бокс')
    }
  }

  return (
    <PartnerLayout
      title="Создать бокс"
      subtitle="Заполните данные предложения для публикации"
      actions={
        <Button variant="secondary" onClick={() => navigate('/partner/boxes')} className="gap-2">
          <ArrowLeft size={16} />
          К списку
        </Button>
      }
    >
      <div className="max-w-3xl">
        <div className="card">
          {locationsStore.loading ? (
            <div className="py-14 flex justify-center">
              <Spinner size={30} />
            </div>
          ) : (
            <BoxForm
              locations={locationsStore.locations}
              onSubmit={handleSubmit}
              isLoading={boxesStore.submitting}
            />
          )}
        </div>
      </div>
    </PartnerLayout>
  )
}

export default observer(CreateBoxPageBase)
