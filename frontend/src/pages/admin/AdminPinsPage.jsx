import { useEffect, useState } from 'react'
import { observer } from 'mobx-react-lite'
import { makeAutoObservable, runInAction } from 'mobx'
import { Plus, Trash2 } from 'lucide-react'
import { toast } from 'sonner'
import AdminLayout from '@/components/admin/layout/AdminLayout'
import Button from '@/components/ui/actions/Button'
import Input from '@/components/ui/form/Input'
import Spinner from '@/components/ui/feedback/Spinner'
import { createAdminLocationPin, deleteAdminLocationPin, listAdminLocationPins } from '@/api/admin'
import { getErrorMessage } from '@/lib/utils'

class AdminPinsStore {
  pins = []
  loading = false
  submitting = false

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
  }

  async load() {
    this.loading = true
    try {
      const data = await listAdminLocationPins()
      runInAction(() => {
        this.pins = data || []
      })
    } finally {
      runInAction(() => {
        this.loading = false
      })
    }
  }

  async create(payload) {
    this.submitting = true
    try {
      await createAdminLocationPin(payload)
      await this.load()
    } finally {
      runInAction(() => {
        this.submitting = false
      })
    }
  }

  async remove(code) {
    await deleteAdminLocationPin(code)
    runInAction(() => {
      this.pins = this.pins.filter((p) => p.code !== code)
    })
  }
}

const store = new AdminPinsStore()

function AdminPinsPageBase() {
  const [code, setCode] = useState('')
  const [nameRu, setNameRu] = useState('')
  const [sortOrder, setSortOrder] = useState('')

  useEffect(() => {
    store.load()
  }, [])

  const handleCreate = async (e) => {
    e.preventDefault()
    if (!code.trim() || !nameRu.trim()) return
    try {
      await store.create({
        code: code.trim(),
        name_ru: nameRu.trim(),
        sort_order: parseInt(sortOrder, 10) || 0,
      })
      setCode('')
      setNameRu('')
      setSortOrder('')
      toast.success('Пин создан')
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleDelete = async (code) => {
    if (!confirm(`Удалить пин "${code}"?`)) return
    try {
      await store.remove(code)
      toast.success('Пин удалён')
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  return (
    <AdminLayout title="Пины заведений">
      <div className="space-y-6">
        {/* Create form */}
        <div className="card p-5">
          <h2 className="text-sm font-semibold text-brand-800 mb-4">Добавить пин</h2>
          <form onSubmit={handleCreate} className="flex flex-wrap gap-3 items-end">
            <div className="flex-1 min-w-36">
              <label className="block text-xs font-medium text-brand-600 mb-1">Код</label>
              <Input
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder="baked_goods"
                required
              />
            </div>
            <div className="flex-1 min-w-36">
              <label className="block text-xs font-medium text-brand-600 mb-1">Название</label>
              <Input
                value={nameRu}
                onChange={(e) => setNameRu(e.target.value)}
                placeholder="Выпечка"
                required
              />
            </div>
            <div className="w-28">
              <label className="block text-xs font-medium text-brand-600 mb-1">Порядок</label>
              <Input
                type="number"
                value={sortOrder}
                onChange={(e) => setSortOrder(e.target.value)}
                placeholder="0"
                min="0"
              />
            </div>
            <Button type="submit" disabled={store.submitting} className="gap-2 shrink-0">
              {store.submitting ? <Spinner size={14} /> : <Plus size={15} />}
              Добавить
            </Button>
          </form>
        </div>

        {/* Pins list */}
        <div className="card divide-y divide-cream-100">
          {store.loading ? (
            <div className="flex justify-center py-10">
              <Spinner size={28} />
            </div>
          ) : store.pins.length === 0 ? (
            <p className="text-center text-brand-500 py-10 text-sm">Пинов нет</p>
          ) : (
            store.pins.map((pin) => (
              <div key={pin.code} className="flex items-center justify-between px-5 py-3.5">
                <div className="flex items-center gap-3">
                  <span className="text-xs font-mono text-brand-400 bg-cream-100 px-2 py-0.5 rounded">{pin.code}</span>
                  <span className="text-sm font-medium text-brand-800">{pin.name_ru}</span>
                  <span className="text-xs text-brand-400">#{pin.sort_order}</span>
                </div>
                <button
                  type="button"
                  onClick={() => handleDelete(pin.code)}
                  className="p-1.5 rounded-lg text-brand-400 hover:text-red-500 hover:bg-red-50 transition-colors"
                  aria-label="Удалить"
                >
                  <Trash2 size={15} />
                </button>
              </div>
            ))
          )}
        </div>
      </div>
    </AdminLayout>
  )
}

export default observer(AdminPinsPageBase)
