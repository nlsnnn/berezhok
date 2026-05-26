import { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { useNavigate } from 'react-router-dom'
import { ArrowLeft } from 'lucide-react'
import { toast } from 'sonner'
import { getErrorMessage } from '@/lib/utils'
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

  // TODO: убрать профиль из locationsStore и юзать только для загрузки локаций, а статус партнера брать из отдельного store

  const partnerStatus = locationsStore.profile?.partner?.status || locationsStore.profile?.status || null
  const canActivateBoxes = partnerStatus !== 'pending_documents'

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
      const created = await boxesStore.create(payload)
      if (formData.status === 'active' && created?.status && created.status !== 'active') {
        toast.info('Бокс сохранён, но не опубликован — заполните юридические данные')
      } else {
        toast.success(formData.status === 'active' ? 'Бокс опубликован' : 'Черновик сохранён')
      }
      navigate('/partner/boxes')
    } catch (error) {
      toast.error(getErrorMessage(error))
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
              mode="create"
              locations={locationsStore.locations}
              onSubmit={handleSubmit}
              isLoading={boxesStore.submitting}
              canActivateBoxes={canActivateBoxes}
            />
          )}
        </div>
      </div>
    </PartnerLayout>
  )
}

export default observer(CreateBoxPageBase)
