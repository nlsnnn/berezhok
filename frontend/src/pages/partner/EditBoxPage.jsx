import { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft } from 'lucide-react'
import { toast } from 'sonner'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import BoxForm from '@/components/partner/boxes/BoxForm'
import Spinner from '@/components/ui/feedback/Spinner'
import Button from '@/components/ui/actions/Button'
import { useStores } from '@/context/StoresContext'

function EditBoxPageBase() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { boxesStore, locationsStore } = useStores()

  useEffect(() => {
    locationsStore.loadProfile()
    boxesStore.loadById(id)
  }, [boxesStore, id, locationsStore])

  const handleSubmit = async (formData) => {
    const payload = {
      name: formData.name,
      description: formData.description,
      original_price: formData.original_price || null,
      discount_price: Number(formData.discount_price),
      quantity_available: Number(formData.quantity_available),
      pickup_time_start: formData.pickup_time_start,
      pickup_time_end: formData.pickup_time_end,
      image_url: formData.image_url || '',
      status: formData.status,
    }

    try {
      await boxesStore.update(id, payload)
      toast.success('Бокс обновлен')
      navigate('/partner/boxes')
    } catch {
      toast.error('Не удалось обновить бокс')
    }
  }

  const isLoading = boxesStore.loading || locationsStore.loading

  return (
    <PartnerLayout
      title="Редактировать бокс"
      subtitle="Обновите параметры предложения"
      actions={
        <Button variant="secondary" onClick={() => navigate('/partner/boxes')} className="gap-2">
          <ArrowLeft size={16} />
          К списку
        </Button>
      }
    >
      <div className="max-w-3xl">
        <div className="card">
          {isLoading ? (
            <div className="py-14 flex justify-center">
              <Spinner size={30} />
            </div>
          ) : boxesStore.current ? (
            <BoxForm
              initialData={boxesStore.current}
              locations={locationsStore.locations}
              onSubmit={handleSubmit}
              isLoading={boxesStore.submitting}
            />
          ) : (
            <div className="py-14 text-center text-red-600">Бокс не найден</div>
          )}
        </div>
      </div>
    </PartnerLayout>
  )
}

export default observer(EditBoxPageBase)
