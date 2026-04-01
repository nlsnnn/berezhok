import { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { useNavigate } from 'react-router-dom'
import { Plus } from 'lucide-react'
import { toast } from 'sonner'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import BoxCard from '@/components/partner/boxes/BoxCard'
import Spinner from '@/components/ui/feedback/Spinner'
import Button from '@/components/ui/actions/Button'
import { useStores } from '@/context/StoresContext'

function BoxesPageBase() {
  const navigate = useNavigate()
  const { boxesStore } = useStores()

  useEffect(() => {
    boxesStore.load()
  }, [boxesStore])

  const handleDelete = async (box) => {
    if (!window.confirm(`Удалить бокс "${box.name}"?`)) return
    try {
      await boxesStore.remove(box.id)
      toast.success('Бокс удален')
    } catch {
      toast.error('Не удалось удалить бокс')
    }
  }

  return (
    <PartnerLayout
      title="Мои боксы"
      subtitle="Управляйте активными и черновыми предложениями"
      actions={
        <Button onClick={() => navigate('/partner/boxes/new')} className="gap-2">
          <Plus size={16} />
          Создать бокс
        </Button>
      }
    >
      {boxesStore.loading && (
        <div className="flex justify-center py-24">
          <Spinner size={34} />
        </div>
      )}

      {boxesStore.error && (
        <div className="card text-center py-12 text-red-600">Не удалось загрузить боксы</div>
      )}

      {!boxesStore.loading && !boxesStore.error && boxesStore.items.length === 0 && (
        <div className="card text-center py-14">
          <p className="text-brand-700 font-medium">Пока нет боксов</p>
          <p className="text-sm text-brand-500 mt-1">Создайте первый сюрприз-бокс для клиентов</p>
          <Button onClick={() => navigate('/partner/boxes/new')} className="mt-5">Создать</Button>
        </div>
      )}

      {!boxesStore.loading && boxesStore.items.length > 0 && (
        <div className="grid sm:grid-cols-2 xl:grid-cols-3 gap-5">
          {boxesStore.items.map((box) => (
            <BoxCard
              key={box.id}
              box={box}
              onEdit={() => navigate(`/partner/boxes/${box.id}/edit`)}
              onDelete={handleDelete}
            />
          ))}
        </div>
      )}
    </PartnerLayout>
  )
}

export default observer(BoxesPageBase)
