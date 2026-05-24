import { useEffect, useState } from 'react'
import { Navigate } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { Edit2, Plus, RefreshCw, ShieldOff } from 'lucide-react'
import AdminBulkActions from '@/components/admin/AdminBulkActions'
import AdminDataState from '@/components/admin/AdminDataState'
import AdminLayout from '@/components/admin/layout/AdminLayout'
import AdminPagination from '@/components/admin/AdminPagination'
import Button from '@/components/ui/actions/Button'
import Input from '@/components/ui/form/Input'
import Select from '@/components/ui/form/Select'
import Modal from '@/components/ui/overlay/Modal'
import StatusBadge from '@/components/ui/feedback/StatusBadge'
import { useStores } from '@/context/StoresContext'
import { canManageAdmins, getAdminRoleLabel } from '@/lib/admin'
import { ADMIN_ACTIVE_STATUS_MAP, ADMIN_EDIT_ROLE_OPTIONS, ADMIN_ROLE_OPTIONS, getStatusMeta } from '@/lib/adminResources'
import { toggleAllPageIds, toggleSelectedId } from '@/lib/adminSelection'
import { formatDateTime, getErrorMessage } from '@/lib/utils'

const emptyForm = {
  email: '',
  password: '',
  name: '',
  role: 'admin',
  is_active: true,
}

function AdminFormModal({ open, admin, loading, onClose, onSubmit }) {
  const isEdit = Boolean(admin)
  const [form, setForm] = useState(() =>
    admin
      ? {
          email: admin.email || '',
          password: '',
          name: admin.name || '',
          role: admin.role || 'admin',
          is_active: admin.is_active !== false,
        }
      : emptyForm
  )

  const update = (name, value) => {
    setForm((current) => ({ ...current, [name]: value }))
  }

  const submit = (event) => {
    event.preventDefault()
    onSubmit(form)
  }

  return (
    <Modal open={open} onClose={onClose} title={isEdit ? 'Редактировать администратора' : 'Новый администратор'} className="max-w-xl">
      <form className="space-y-4" onSubmit={submit}>
        <label className="block space-y-2">
          <span className="text-sm font-medium text-brand-700">Email</span>
          <Input type="email" value={form.email} onChange={(event) => update('email', event.target.value)} required disabled={isEdit} />
        </label>
        {!isEdit && (
          <label className="block space-y-2">
            <span className="text-sm font-medium text-brand-700">Пароль</span>
            <Input type="password" value={form.password} onChange={(event) => update('password', event.target.value)} minLength={8} required />
          </label>
        )}
        <label className="block space-y-2">
          <span className="text-sm font-medium text-brand-700">Имя</span>
          <Input value={form.name} onChange={(event) => update('name', event.target.value)} required />
        </label>
        <label className="block space-y-2">
          <span className="text-sm font-medium text-brand-700">Роль</span>
          <Select value={form.role} onChange={(event) => update('role', event.target.value)}>
            {(isEdit ? ADMIN_EDIT_ROLE_OPTIONS : ADMIN_ROLE_OPTIONS).map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </Select>
        </label>
        {isEdit && (
          <label className="flex items-center gap-3 rounded-xl border border-cream-200 bg-cream-50 p-3 text-sm text-brand-700">
            <input type="checkbox" checked={form.is_active} onChange={(event) => update('is_active', event.target.checked)} />
            Активен
          </label>
        )}
        <div className="flex justify-end gap-3 border-t border-cream-200 pt-4">
          <Button type="button" variant="secondary" onClick={onClose}>Отмена</Button>
          <Button type="submit" disabled={loading}>{isEdit ? 'Сохранить' : 'Создать'}</Button>
        </div>
      </form>
    </Modal>
  )
}

function AdminAdminsPageBase() {
  const { adminAdminsStore, adminAuthStore } = useStores()
  const [modalOpen, setModalOpen] = useState(false)
  const [selected, setSelected] = useState(null)
  const [selectedIds, setSelectedIds] = useState([])

  useEffect(() => {
    if (canManageAdmins(adminAuthStore.user)) {
      adminAdminsStore.load().catch(() => {})
    }
  }, [adminAdminsStore, adminAuthStore.user])

  if (!canManageAdmins(adminAuthStore.user)) {
    return <Navigate to="/admin/applications" replace />
  }

  const openCreate = () => {
    setSelected(null)
    setModalOpen(true)
  }

  const openEdit = (admin) => {
    setSelected(admin)
    setModalOpen(true)
  }

  const closeModal = () => {
    setSelected(null)
    setModalOpen(false)
  }

  const submit = async (form) => {
    try {
      if (selected) {
        await adminAdminsStore.update(selected.id || selected.user_id, {
          name: form.name,
          role: form.role,
          is_active: form.is_active,
        })
        toast.success('Администратор обновлён')
      } else {
        await adminAdminsStore.create({
          email: form.email,
          password: form.password,
          name: form.name,
          role: form.role,
        })
        toast.success('Администратор создан')
      }
      closeModal()
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const deactivate = async (admin) => {
    if (!window.confirm('Деактивировать администратора?')) return

    try {
      await adminAdminsStore.deactivate(admin.id || admin.user_id)
      toast.success('Администратор деактивирован')
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const deactivateSelected = async (ids) => {
    if (!window.confirm(`Деактивировать выбранных администраторов: ${ids.length}?`)) return

    try {
      await adminAdminsStore.deactivateMany(ids)
      toast.success('Выбранные администраторы деактивированы')
      setSelectedIds([])
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const pageIds = adminAdminsStore.items.map((admin) => admin.id || admin.user_id).filter(Boolean)
  const selectedPageIds = selectedIds.filter((id) => pageIds.includes(id))
  const allPageSelected = pageIds.length > 0 && pageIds.every((id) => selectedPageIds.includes(id))

  return (
    <AdminLayout
      title="Администраторы"
      subtitle="Управление доступом к админ-панели"
      actions={
        <>
          <Button variant="secondary" onClick={() => adminAdminsStore.load()} disabled={adminAdminsStore.loading}>
            <RefreshCw size={16} />
            Обновить
          </Button>
          <Button onClick={openCreate}>
            <Plus size={16} />
            Добавить
          </Button>
        </>
      }
    >
      <AdminDataState
        loading={adminAdminsStore.loading}
        error={adminAdminsStore.error}
        empty={adminAdminsStore.items.length === 0}
        emptyText="Администраторов пока нет"
        onRetry={() => adminAdminsStore.load()}
      >
        <section className="space-y-4">
          <AdminBulkActions
            selectedIds={selectedPageIds}
            pageIds={pageIds}
            allPageSelected={allPageSelected}
            onToggleAll={() => setSelectedIds((current) => toggleAllPageIds(current, pageIds))}
            onClear={() => setSelectedIds((current) => current.filter((id) => !pageIds.includes(id)))}
            actions={[
              {
                label: 'Деактивировать',
                icon: <ShieldOff size={16} />,
                variant: 'danger',
                onClick: deactivateSelected,
              },
            ]}
            loading={adminAdminsStore.loading || adminAdminsStore.actionLoading}
          />

          {adminAdminsStore.items.map((admin) => {
            const adminId = admin.id || admin.user_id
            const status = admin.is_active === false ? 'inactive' : 'active'
            const statusMeta = getStatusMeta(status, ADMIN_ACTIVE_STATUS_MAP)

            return (
              <article key={adminId} className="card">
                <div className="flex flex-wrap items-start justify-between gap-4">
                  <label className="flex items-center pt-1">
                    <input
                      type="checkbox"
                      className="h-4 w-4 rounded border-cream-300 text-brand-600"
                      checked={selectedPageIds.includes(adminId)}
                      onChange={() => setSelectedIds((current) => toggleSelectedId(current, adminId))}
                      aria-label="Выбрать администратора"
                    />
                  </label>
                  <div>
                    <div className="mb-2 flex flex-wrap items-center gap-3">
                      <StatusBadge status={status} customLabel={statusMeta.label} customColor={statusMeta.color} />
                      <span className="text-xs text-cream-500">{getAdminRoleLabel(admin.role)}</span>
                    </div>
                    <h2 className="text-base font-semibold text-brand-900">{admin.name || 'Администратор'}</h2>
                    <p className="mt-1 text-sm text-brand-600">{admin.email}</p>
                    <p className="mt-2 text-xs text-brand-500">Создан {formatDateTime(admin.created_at)}</p>
                  </div>
                  <div className="flex flex-wrap gap-2">
                    <Button variant="secondary" size="sm" onClick={() => openEdit(admin)}>
                      <Edit2 size={16} />
                      Изменить
                    </Button>
                    {admin.is_active !== false && (
                      <Button variant="danger" size="sm" onClick={() => deactivate(admin)} disabled={adminAdminsStore.actionLoading}>
                        <ShieldOff size={16} />
                        Деактивировать
                      </Button>
                    )}
                  </div>
                </div>
              </article>
            )
          })}
          <AdminPagination pagination={adminAdminsStore.pagination} loading={adminAdminsStore.loading} onPrev={() => adminAdminsStore.prevPage()} onNext={() => adminAdminsStore.nextPage()} />
        </section>
      </AdminDataState>

      <AdminFormModal
        key={selected?.id || selected?.user_id || 'new-admin'}
        open={modalOpen}
        admin={selected}
        loading={adminAdminsStore.actionLoading}
        onClose={closeModal}
        onSubmit={submit}
      />
    </AdminLayout>
  )
}

export default observer(AdminAdminsPageBase)
