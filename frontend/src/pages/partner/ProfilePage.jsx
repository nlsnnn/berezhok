import { observer } from 'mobx-react-lite'
import { UserCircle } from 'lucide-react'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'

function ProfilePageBase() {
  return (
    <PartnerLayout
      title="Профиль"
      subtitle="Настройки аккаунта и информация о партнёре"
    >
      <div className="card flex flex-col items-center justify-center py-16 text-center">
        <div className="w-14 h-14 rounded-2xl bg-brand-100 flex items-center justify-center mb-4">
          <UserCircle size={24} className="text-brand-600" />
        </div>
        <h3 className="text-lg font-semibold text-brand-900 mb-2">Раздел в разработке</h3>
        <p className="text-sm text-brand-600 max-w-sm">
          Здесь вы сможете редактировать данные профиля, менять пароль и управлять настройками аккаунта.
        </p>
      </div>
    </PartnerLayout>
  )
}

export default observer(ProfilePageBase)
