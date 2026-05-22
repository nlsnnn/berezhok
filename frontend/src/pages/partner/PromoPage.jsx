import { useState } from 'react'
import { Plus, Gift, Info, MoreVertical, Calendar } from 'lucide-react'
import Button from '@/components/ui/actions/Button'
import Input from '@/components/ui/form/Input'
import Select from '@/components/ui/form/Select'
import Modal from '@/components/ui/overlay/Modal'

// Моковые данные промокодов для демонстрации
const mockPromos = [
  {
    id: 1,
    code: 'NEWBIE15',
    type: 'percent',
    value: 15,
    appliesTo: 'all',
    restrictions: { newUsers: true, minOrder: 1000 },
    stats: { used: 42 },
    active: true,
  },
  {
    id: 2,
    code: 'MINUS50',
    type: 'fixed',
    value: 50,
    appliesTo: 'box',
    boxName: 'Сладкий бокс из пекарни',
    restrictions: { oncePerUser: true, dateLimit: '2026-06-01' },
    stats: { used: 120 },
    active: true,
  },
  {
    id: 3,
    code: 'FIRSTLUCKY',
    type: 'free',
    value: 100,
    appliesTo: 'all',
    restrictions: { totalLimit: 1 },
    stats: { used: 1 },
    active: false,
  },
]

export default function PromoPage() {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [activeTab, setActiveTab] = useState('active')

  // Стэйт для формы создания (как макет)
  const [promoType, setPromoType] = useState('percent')
  const [appliesTo, setAppliesTo] = useState('all')

  const filteredPromos = mockPromos.filter((p) =>
    activeTab === 'active' ? p.active : !p.active
  )

  return (
    <div className="space-y-6 max-w-5xl mx-auto pb-8 relative">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold text-brand-900">Промокоды</h1>
        <Button onClick={() => setIsModalOpen(true)}>
          <Plus size={18} className="mr-2" />
          Создать
        </Button>
      </div>

      <div className="flex space-x-2 border-b border-cream-200">
        <button
          className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
            activeTab === 'active'
              ? 'border-brand-600 text-brand-600'
              : 'border-transparent text-gray-500 hover:text-gray-700'
          }`}
          onClick={() => setActiveTab('active')}
        >
          Активные
        </button>
        <button
          className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
            activeTab === 'inactive'
              ? 'border-brand-600 text-brand-600'
              : 'border-transparent text-gray-500 hover:text-gray-700'
          }`}
          onClick={() => setActiveTab('inactive')}
        >
          Неактивные
        </button>
      </div>

      {filteredPromos.length === 0 ? (
        <div className="text-center py-12 bg-white rounded-xl shadow-sm border border-cream-200">
          <Gift size={48} className="mx-auto text-cream-300 mb-4" />
          <h3 className="text-lg font-medium text-brand-900">Нет промокодов</h3>
          <p className="text-gray-500 mt-1">Здесь пока ничего нет</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {filteredPromos.map((promo) => (
            <div
              key={promo.id}
              className="bg-white rounded-xl shadow-sm border border-cream-200 p-5 relative"
            >
              <div className="flex justify-between items-start mb-3">
                <div className="flex-1">
                  <span className="inline-block px-2.5 py-1 text-xs font-semibold uppercase tracking-wide bg-brand-50 text-brand-700 rounded-md">
                    {promo.code}
                  </span>
                </div>
                <button className="text-gray-400 hover:text-gray-600">
                  <MoreVertical size={20} />
                </button>
              </div>

              <div className="mb-4">
                <div className="text-2xl font-bold text-gray-900">
                  {promo.type === 'percent' && `-${promo.value}%`}
                  {promo.type === 'fixed' && `-${promo.value}₽`}
                  {promo.type === 'free' && `Бесплатно`}
                </div>
                <p className="text-sm text-gray-500 mt-1">
                  {promo.appliesTo === 'all'
                    ? 'На все боксы заведения'
                    : `Только на бокс: ${promo.boxName}`}
                </p>
              </div>

              <div className="space-y-2 mb-5">
                {promo.restrictions?.newUsers && (
                  <div className="flex items-center text-xs text-gray-600">
                    <Info size={14} className="mr-1.5 text-blue-500" />
                    Только для новых пользователей
                  </div>
                )}
                {promo.restrictions?.oncePerUser && (
                  <div className="flex items-center text-xs text-gray-600">
                    <Info size={14} className="mr-1.5 text-blue-500" />
                    Один раз на пользователя
                  </div>
                )}
                {promo.restrictions?.totalLimit && (
                  <div className="flex items-center text-xs text-gray-600">
                    <Info size={14} className="mr-1.5 text-orange-500" />
                    Ограничено до {promo.restrictions.totalLimit} использований
                  </div>
                )}
                {promo.restrictions?.dateLimit && (
                  <div className="flex items-center text-xs text-gray-600">
                    <Calendar size={14} className="mr-1.5 text-red-500" />
                    Действует до {promo.restrictions.dateLimit}
                  </div>
                )}
              </div>

              <div className="border-t border-cream-100 pt-4 mt-auto">
                <div className="flex justify-between items-center text-sm">
                  <span className="text-gray-500">Использовано:</span>
                  <span className="font-semibold text-gray-900">
                    {promo.stats.used} раз
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      <Modal
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Создание промокода"
      >
        <div className="space-y-5">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Код (будет сгенерирован автоматически, если оставить пустым)
            </label>
            <Input placeholder="Например, HAPPYNEWYEAR" />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Тип скидки</label>
              <Select value={promoType} onChange={(e) => setPromoType(e.target.value)}>
                <option value="percent">Процент (%)</option>
                <option value="fixed">Сумма (₽)</option>
                <option value="free">Бесплатный бокс (100%)</option>
              </Select>
            </div>
            {promoType !== 'free' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Размер</label>
                <Input type="number" placeholder="15" />
              </div>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Применение</label>
            <Select value={appliesTo} onChange={(e) => setAppliesTo(e.target.value)}>
              <option value="all">На все боксы заведения</option>
              <option value="box">На конкретный бокс</option>
            </Select>
          </div>

          {appliesTo === 'box' && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Выберите бокс</label>
              <Select>
                <option value="">-- Выберите из списка --</option>
                <option value="1">Сладкий бокс из пекарни</option>
                <option value="2">Обед комбо</option>
              </Select>
            </div>
          )}

          <div className="border-t border-cream-200 mt-6 pt-4 space-y-4">
            <h4 className="font-medium text-brand-900">Ограничения</h4>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Всего использований</label>
                <Input type="number" placeholder="Не ограничено" />
                <p className="text-xs text-gray-500 mt-1">Один раз всего = 1</p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Действует до</label>
                <Input type="date" />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Минимальный чек (₽)</label>
                <Input type="number" placeholder="0" />
              </div>
            </div>

            <div className="space-y-3 pt-2">
              <label className="flex items-center space-x-3 cursor-pointer">
                <input
                  type="checkbox"
                  className="w-4 h-4 text-brand-600 border-gray-300 rounded focus:ring-brand-500"
                />
                <span className="text-sm text-gray-700">Один раз на пользователя</span>
              </label>

              <label className="flex items-center space-x-3 cursor-pointer">
                <input
                  type="checkbox"
                  className="w-4 h-4 text-brand-600 border-gray-300 rounded focus:ring-brand-500"
                />
                <span className="text-sm text-gray-700">Только для новых пользователей (нет завершенных заказов)</span>
              </label>
            </div>
          </div>

          <div className="flex justify-end pt-5 space-x-3">
            <Button variant="secondary" onClick={() => setIsModalOpen(false)}>
              Отмена
            </Button>
            <Button onClick={() => setIsModalOpen(false)}>
              Создать
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  )
}
