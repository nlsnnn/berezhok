import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { listBoxes, deleteBox } from '@/api/partner'
import { Plus, Package } from 'lucide-react'
import { toast } from 'sonner'
import PartnerNav from '@/components/PartnerNav'
import Spinner from '@/components/ui/Spinner'
import Button from '@/components/ui/Button'
import BoxCard from '@/components/ui/BoxCard'

export default function BoxesPage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data: boxes, isLoading, isError } = useQuery({
    queryKey: ['partner', 'boxes'],
    queryFn: listBoxes,
  })

  const deleteMutation = useMutation({
    mutationFn: deleteBox,
    onSuccess: () => {
      queryClient.invalidateQueries(['partner', 'boxes'])
      toast.success('Бокс удалён')
    },
    onError: (error) => {
      console.error('Delete error:', error)
      toast.error('Не удалось удалить бокс')
    },
  })

  const handleEdit = (box) => {
    navigate(`/partner/boxes/${box.id}/edit`)
  }

  const handleDelete = (box) => {
    if (window.confirm(`Удалить бокс "${box.name}"?`)) {
      deleteMutation.mutate(box.id)
    }
  }

  return (
    <div className="min-h-screen flex flex-col bg-cream-50">
      <PartnerNav />

      <main className="flex-1 max-w-7xl mx-auto w-full px-4 sm:px-6 py-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-3xl font-bold text-brand-900">Мои боксы</h1>
            <p className="text-brand-600 mt-1">
              Управляйте своими предложениями для клиентов
            </p>
          </div>
          <Button
            onClick={() => navigate('/partner/boxes/new')}
            className="gap-2"
          >
            <Plus size={18} />
            Создать бокс
          </Button>
        </div>

        {/* Loading */}
        {isLoading && (
          <div className="flex justify-center py-20">
            <Spinner size={32} />
          </div>
        )}

        {/* Error */}
        {isError && (
          <div className="card text-center py-12 text-red-500">
            Не удалось загрузить боксы
          </div>
        )}

        {/* Empty state */}
        {boxes && boxes.length === 0 && (
          <div className="card text-center py-16">
            <div className="w-20 h-20 rounded-full bg-cream-200 flex items-center justify-center mx-auto mb-4">
              <Package size={40} className="text-cream-400" />
            </div>
            <h3 className="font-semibold text-brand-700 text-lg mb-2">
              Нет боксов
            </h3>
            <p className="text-brand-500 mb-6">
              Создайте свой первый бокс-сюрприз для клиентов
            </p>
            <Button
              onClick={() => navigate('/partner/boxes/new')}
              className="gap-2"
            >
              <Plus size={18} />
              Создать бокс
            </Button>
          </div>
        )}

        {/* Boxes grid */}
        {boxes && boxes.length > 0 && (
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {boxes.map((box) => (
              <BoxCard
                key={box.id}
                box={box}
                onEdit={handleEdit}
                onDelete={handleDelete}
              />
            ))}
          </div>
        )}
      </main>
    </div>
  )
}
